package imp


import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	pongo2 "github.com/flosch/pongo2"
	"os/exec"
	"fmt"
	"os"
	"io/ioutil"
	"net"
	"errors"
	"bytes"
	"strings"
	config "installer/config"
	yaml "gopkg.in/yaml.v2"
    jsonyaml "github.com/ghodss/yaml"
	v1 "k8s.io/api/core/v1"
)

func Run(name string,configPath string,hostPkiPath string) {


    //读取安装配置文件
	installConfig :=loadConfig(configPath)
	//取得主机名
	nodeHostName,err:= os.Hostname()
	check(err)
	nodeName := name
	var nodeIpAddress string
	if currentNode, ok := installConfig.Nodes[nodeName]; ok {
		if len(currentNode.Ipv4Address)!=0{
			nodeIpAddress = currentNode.Ipv4Address
		}else{
			nodeIpAddress,err = getIPV4()
			check(err)
			currentNode.Ipv4Address = nodeIpAddress
			installConfig.Nodes[nodeName] = currentNode
		}
	}else{
		panic(fmt.Sprintf("Can not find node: %s in install config file",name))
	}
	var etcdPeers []string
	var etcdEndpoints []string
	for _,nname:=range(installConfig.Etcds){
		if currentNode, ok := installConfig.Nodes[nname]; ok {
			if len(currentNode.Ipv4Address)!=0{
				etcdPeers = append(etcdPeers,fmt.Sprintf("%s=https://%s:2380",nname,currentNode.Ipv4Address))
				etcdEndpoints = append(etcdEndpoints,fmt.Sprintf("https://%s:2379",currentNode.Ipv4Address))
			}else{
				panic(fmt.Sprintf("node:%s nodeIpAddress is blank",name))
			}
		}else{
			panic(fmt.Sprintf("Can not find node:%s in install config nodes",name))
		}

	}

	pkiPath :="/etc/kubernetes/pki"
	kubeadmConfigPath :="/etc/kubernetes/kubeadm"
	certTplPath := "/root/template/cert"
	kubeadmTplPath := "/root/template/config"


	var kubeadmContext = pongo2.Context{
		"hostPkiPath":hostPkiPath,
		"pkiPath": pkiPath,
		"nodeName": nodeName,
		"nodeHostName": nodeHostName,
		"nodeIpAddress": nodeIpAddress,
		"dnsDomain": "cluster.local",
		"serviceSubnet":"10.233.0.0/18",
		"podSubnet":"10.233.64.0/18",
		"kubernetesVersion":"1.12.3",
		"kubeImageRepo":"registry.cn-hangzhou.aliyuncs.com/k8s_mirror",
		"token": "vm83ja.l2ri182gvmfx2168",
		"etcdEndpoints":etcdEndpoints,
		"authorizationModes":[]string{"Node","RBAC"},
		"admissionControl":strings.Join([]string{"Initializers","NamespaceLifecycle","LimitRanger","ServiceAccount","DefaultStorageClass","GenericAdmissionWebhook","ResourceQuota"},","),
		"etcdClusterPeers" : strings.Join(etcdPeers,","),
		"etcdState" : "new",
		"etcdToken" : "choerodon-install-etcd",
		}

	parserDir(certTplPath,pkiPath+"/config",kubeadmContext)


	parserDir(kubeadmTplPath+"/master",kubeadmConfigPath+"/master",kubeadmContext)
	//生成CA证书
	execShellFile(pkiPath+"/config/"+"ca_gen_cert.sh")


	//生成ETCD证书
	execShellFile(pkiPath+"/config/"+"etcd_gen_cert_client.sh")
	execShellFile(pkiPath+"/config/"+"etcd_gen_cert_server.sh")

	parserDir("/root/template/docker/etcd","/etc/kubernetes/docker/etcd",kubeadmContext)
	execShellFile("/etc/kubernetes/docker/etcd/run.sh")


	//生成front-proxy-ca证书
	execShellFile(pkiPath+"/config/"+"k8s_gen_cert_front_proxy_ca.sh")

	//生成api-server证书
	execShellFile(pkiPath+"/config/"+"k8s_gen_cert_apiserver.sh")

	//生成sa证书
	execShellFile(pkiPath+"/config/"+"k8s_gen_cert_sa.sh")

	//生成 pod yml
	execShellFile(pkiPath+"/config/"+"kubeadm.sh")


}

func loadConfig(configPath string) *config.Config{
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatal(fmt.Sprintf("Can not find install config file at %s",configPath))
		panic(fmt.Sprintf("Can not file install config file at %s",configPath))
	}
	installConfig := &config.Config{}
	buffer, err := ioutil.ReadFile(configPath)
	check(err)
	err = yaml.Unmarshal(buffer, &installConfig)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	return installConfig
}

func check(e error) {
	if e != nil {
		log.Fatal(e)
		panic(e)
	}
}

func execShellFile(filePath string) []byte {
	log.Info("Execute file:"+filePath)
	cmd := exec.Command("/bin/sh", filePath)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		log.Fatal(fmt.Sprint(err) + ": " + stderr.String())
		panic(err)
	}

	log.Info(fmt.Sprintf("%s",out))
	return out.Bytes()

}

func getIPV4() (string,error){
	addrs, err := net.InterfaceAddrs()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for _, address := range addrs {

		// check the address type and if it is not a loopback the display it
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				fmt.Println(ipnet.IP.String())
				return ipnet.IP.String(),nil
			}

		}
	}
	return "", errors.New("ipv4: can not find a ip address")


}

func parserDir(templateDir string,toDir string,context pongo2.Context){
	if _, err := os.Stat(toDir); os.IsNotExist(err) {
		e := os.MkdirAll(toDir, 0700)
		check(e)
	}
	log.Info(fmt.Sprintf("parser directory %s",templateDir))
	dir_list, err := ioutil.ReadDir(templateDir)
	if err != nil {
		fmt.Println("read dir error")
		return
	}
	for _,v := range dir_list {
		if !v.IsDir(){
			parserTemplate(templateDir+"/"+v.Name(),toDir+"/"+v.Name(),context)
		}
	}
}

func parserTemplate(templatePath string,toPath string,context pongo2.Context){
	log.Info(fmt.Sprintf("parser template %s",templatePath))
	deleteFileIfExists(toPath)
	var tpl = pongo2.Must(pongo2.FromFile(templatePath))
	f, err := os.Create(toPath)
	defer f.Close()
	check(err)
	err = tpl.ExecuteWriter(context,f)
	check(err)
}

func deleteFileIfExists(filePath string) {
	os.Remove(filePath)
}

func podToDocker(filePath string){
	fileBytes,err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Println("read dir error")
		return
	}
	jsonBytes, err := jsonyaml.YAMLToJSON(fileBytes)

	if err != nil {
	return v1.Job{}, err
	}
	// unmarshal the json into the kube struct
	var job = v1.Job{}
	err = json.Unmarshal(jsonBytes, &job)
	if err != nil {
	return v1.Job{}, err
	}
}


// yamlBytes contains a []byte of my yaml job spec
// convert the yaml to json
//jsonBytes, err := yaml.YAMLToJSON(yamlBytes)
//if err != nil {
//return v1.Job{}, err
//}
//// unmarshal the json into the kube struct
//var job = v1.Job{}
//err = json.Unmarshal(jsonBytes, &job)
//if err != nil {
//return v1.Job{}, err
//}
// job now contains the for the job that was in the YAML!