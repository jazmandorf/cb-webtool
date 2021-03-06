package service

import (
	"encoding/json"
	"fmt"
	// "io"
	"log"
	"net/http"

	// "os"
	// model "github.com/cloud-barista/cb-webtool/src/model"
	"github.com/cloud-barista/cb-webtool/src/model"
	// spider "github.com/cloud-barista/cb-webtool/src/model/spider"
	"github.com/cloud-barista/cb-webtool/src/model/ladybug"
	// "github.com/cloud-barista/cb-webtool/src/model/tumblebug"

	util "github.com/cloud-barista/cb-webtool/src/util"
)

// Health Check
func GetHealthy() model.WebStatus {
	var originalUrl = "/healthy"
	urlParam := util.MappingUrlParameter(originalUrl, nil)
	url := util.LADYBUG + urlParam
	resp, err := util.CommonHttp(url, nil, http.MethodGet)
	// resp, err := util.CommonHttpWithoutParam(url, http.MethodGet)

	if err != nil {
		fmt.Println(err)
		return model.WebStatus{StatusCode: 500, Message: err.Error()}
	}

	respBody := resp.Body
	respStatus := resp.StatusCode

	result := ""
	json.NewDecoder(respBody).Decode(&result)
	log.Println(result)

	return model.WebStatus{StatusCode: respStatus, Message: result}
}

// Cluster 목록 조회
func GetClusterList(nameSpaceID string) ([]ladybug.ClusterInfo, model.WebStatus) {
	var originalUrl = "/ns/{namespace}/clusters"

	var paramMapper = make(map[string]string)
	paramMapper["{namespace}"] = nameSpaceID
	urlParam := util.MappingUrlParameter(originalUrl, paramMapper)

	url := util.LADYBUG + urlParam
	// url := util.LADYBUG + "/ns/" + nameSpaceID + "/clusters"
	resp, err := util.CommonHttp(url, nil, http.MethodGet)
	// resp, err := util.CommonHttpWithoutParam(url, http.MethodGet)

	if err != nil {
		fmt.Println(err)
		return nil, model.WebStatus{StatusCode: 500, Message: err.Error()}
	}

	respBody := resp.Body
	respStatus := resp.StatusCode

	clusterList := map[string][]ladybug.ClusterInfo{}
	json.NewDecoder(respBody).Decode(&clusterList)
	fmt.Println(clusterList["items"])
	log.Println(respBody)
	// util.DisplayResponse(resp) // 수신내용 확인

	return clusterList["items"], model.WebStatus{StatusCode: respStatus}
}

// 특정 Cluster 조회
func GetClusterData(nameSpaceID string, cluster string) (*ladybug.ClusterInfo, model.WebStatus) {
	var originalUrl = "/ns/{namespace}/clusters/{cluster}"

	var paramMapper = make(map[string]string)
	paramMapper["{namespace}"] = nameSpaceID
	paramMapper["{cluster}"] = cluster
	urlParam := util.MappingUrlParameter(originalUrl, paramMapper)

	url := util.LADYBUG + urlParam

	// resp, err := util.CommonHttp(url, nil, http.MethodGet)
	resp, err := util.CommonHttpWithoutParam(url, http.MethodGet)

	// defer body.Close()
	clusterInfo := ladybug.ClusterInfo{}
	if err != nil {
		fmt.Println(err)
		return &clusterInfo, model.WebStatus{StatusCode: 500, Message: err.Error()}
	}
	// util.DisplayResponse(resp) // 수신내용 확인

	respBody := resp.Body
	respStatus := resp.StatusCode

	json.NewDecoder(respBody).Decode(&clusterInfo)
	fmt.Println(clusterInfo)

	return &clusterInfo, model.WebStatus{StatusCode: respStatus}
}

// Cluster 생성
func RegCluster(nameSpaceID string, clusterReq *ladybug.ClusterReq) (*ladybug.ClusterInfo, model.WebStatus) {

	var originalUrl = "/ns/{namespace}/clusters"

	var paramMapper = make(map[string]string)
	paramMapper["{namespace}"] = nameSpaceID
	urlParam := util.MappingUrlParameter(originalUrl, paramMapper)
	url := util.LADYBUG + urlParam

	pbytes, _ := json.Marshal(clusterReq)
	resp, err := util.CommonHttp(url, pbytes, http.MethodPost)

	returnClusterInfo := ladybug.ClusterInfo{}
	returnStatus := model.WebStatus{}

	respBody := resp.Body
	respStatus := resp.StatusCode

	if err != nil {
		fmt.Println(err)
		return &returnClusterInfo, model.WebStatus{StatusCode: 500, Message: err.Error()}
	}

	if respStatus != 200 && respStatus != 201 { // 호출은 정상이나, 가져온 결과값이 200, 201아닌 경우 message에 담겨있는 것을 WebStatus에 set
		errorInfo := model.ErrorInfo{}
		json.NewDecoder(respBody).Decode(&errorInfo)
		fmt.Println("respStatus != 200 reason ", errorInfo)
		returnStatus.Message = errorInfo.Message
	} else {
		json.NewDecoder(respBody).Decode(&returnClusterInfo)
		fmt.Println(returnClusterInfo)
	}
	returnStatus.StatusCode = respStatus

	return &returnClusterInfo, returnStatus
}

// Cluster 삭제
func DelCluster(nameSpaceID string, cluster string) (*ladybug.StatusInfo, model.WebStatus) {
	var originalUrl = "/ns/{namespace}/clusters/{cluster}"

	var paramMapper = make(map[string]string)
	paramMapper["{namespace}"] = nameSpaceID
	paramMapper["{cluster}"] = cluster
	urlParam := util.MappingUrlParameter(originalUrl, paramMapper)
	url := util.LADYBUG + urlParam

	if cluster == "" {
		return nil, model.WebStatus{StatusCode: 500, Message: "cluster is required"}
	}

	// 경로안에 parameter가 있어 추가 param없이 호출 함.
	resp, err := util.CommonHttp(url, nil, http.MethodDelete)
	statusInfo := ladybug.StatusInfo{}
	if err != nil {
		fmt.Println(err)
		return &statusInfo, model.WebStatus{StatusCode: 500, Message: err.Error()}
	}

	if err != nil {
		fmt.Println(err)
		return &statusInfo, model.WebStatus{StatusCode: 500, Message: err.Error()}
	}
	// util.DisplayResponse(resp) // 수신내용 확인

	respBody := resp.Body
	respStatus := resp.StatusCode

	json.NewDecoder(respBody).Decode(&statusInfo)
	fmt.Println(statusInfo)

	return &statusInfo, model.WebStatus{StatusCode: respStatus}
}

////////

// Node 목록 조회
func GetNodeList(nameSpaceID string, cluster string) (ladybug.NodeList, model.WebStatus) {
	var originalUrl = "/ns/{namespace}/clusters/{cluster}/nodes"

	var paramMapper = make(map[string]string)
	paramMapper["{namespace}"] = nameSpaceID
	paramMapper["{cluster}"] = cluster
	urlParam := util.MappingUrlParameter(originalUrl, paramMapper)

	url := util.LADYBUG + urlParam

	resp, err := util.CommonHttp(url, nil, http.MethodGet)
	// resp, err := util.CommonHttpWithoutParam(url, http.MethodGet)

	nodeList := ladybug.NodeList{} // 이름은 List이나 1개의 객체임
	if err != nil {
		fmt.Println(err)
		return nodeList, model.WebStatus{StatusCode: 500, Message: err.Error()}
	}

	respBody := resp.Body
	respStatus := resp.StatusCode

	json.NewDecoder(respBody).Decode(&nodeList)
	fmt.Println(nodeList)
	log.Println(respBody)
	// util.DisplayResponse(resp) // 수신내용 확인

	return nodeList, model.WebStatus{StatusCode: respStatus}
}

// 특정 Cluster 조회
func GetNodeData(nameSpaceID string, cluster string, node string) (*ladybug.NodeInfo, model.WebStatus) {
	var originalUrl = "/ns/{namespace}/clusters/{cluster}/nodes/{node}"

	var paramMapper = make(map[string]string)
	paramMapper["{namespace}"] = nameSpaceID
	paramMapper["{cluster}"] = cluster
	paramMapper["{node}"] = node
	urlParam := util.MappingUrlParameter(originalUrl, paramMapper)

	url := util.LADYBUG + urlParam

	// resp, err := util.CommonHttp(url, nil, http.MethodGet)
	resp, err := util.CommonHttpWithoutParam(url, http.MethodGet)

	// defer body.Close()
	nodeInfo := ladybug.NodeInfo{}
	if err != nil {
		fmt.Println(err)
		return &nodeInfo, model.WebStatus{StatusCode: 500, Message: err.Error()}
	}
	// util.DisplayResponse(resp) // 수신내용 확인

	respBody := resp.Body
	respStatus := resp.StatusCode

	json.NewDecoder(respBody).Decode(&nodeInfo)
	fmt.Println(nodeInfo)

	return &nodeInfo, model.WebStatus{StatusCode: respStatus}
}

// Node 생성
func RegNode(nameSpaceID string, cluster string, nodeReq *ladybug.NodeReq) (*ladybug.NodeInfo, model.WebStatus) {

	var originalUrl = "/ns/{namespace}/clusters/{cluster}/nodes"

	var paramMapper = make(map[string]string)
	paramMapper["{namespace}"] = nameSpaceID
	paramMapper["{cluster}"] = cluster
	urlParam := util.MappingUrlParameter(originalUrl, paramMapper)
	url := util.LADYBUG + urlParam

	pbytes, _ := json.Marshal(nodeReq)
	resp, err := util.CommonHttp(url, pbytes, http.MethodPost)

	returnNodeInfo := ladybug.NodeInfo{}
	returnStatus := model.WebStatus{}

	respBody := resp.Body
	respStatus := resp.StatusCode

	if err != nil {
		fmt.Println(err)
		return &returnNodeInfo, model.WebStatus{StatusCode: 500, Message: err.Error()}
	}

	if respStatus != 200 && respStatus != 201 { // 호출은 정상이나, 가져온 결과값이 200, 201아닌 경우 message에 담겨있는 것을 WebStatus에 set
		errorInfo := model.ErrorInfo{}
		json.NewDecoder(respBody).Decode(&errorInfo)
		fmt.Println("respStatus != 200 reason ", errorInfo)
		returnStatus.Message = errorInfo.Message
	} else {
		json.NewDecoder(respBody).Decode(&returnNodeInfo)
		fmt.Println(returnNodeInfo)
	}
	returnStatus.StatusCode = respStatus

	return &returnNodeInfo, returnStatus
}

// Node 삭제
func DelNode(nameSpaceID string, cluster string, node string) (*ladybug.StatusInfo, model.WebStatus) {
	var originalUrl = "/ns/{namespace}/clusters/{cluster}/nodes/{node}"

	var paramMapper = make(map[string]string)
	paramMapper["{namespace}"] = nameSpaceID
	paramMapper["{cluster}"] = cluster
	paramMapper["{node}"] = node
	urlParam := util.MappingUrlParameter(originalUrl, paramMapper)
	url := util.LADYBUG + urlParam

	if cluster == "" {
		return nil, model.WebStatus{StatusCode: 500, Message: "cluster is required"}
	}
	if node == "" {
		return nil, model.WebStatus{StatusCode: 500, Message: "node is required"}
	}

	// 경로안에 parameter가 있어 추가 param없이 호출 함.
	resp, err := util.CommonHttp(url, nil, http.MethodDelete)
	if err != nil {
		fmt.Println(err)
		return nil, model.WebStatus{StatusCode: 500, Message: err.Error()}
	}

	statusInfo := ladybug.StatusInfo{}
	if err != nil {
		fmt.Println(err)
		return &statusInfo, model.WebStatus{StatusCode: 500, Message: err.Error()}
	}
	// util.DisplayResponse(resp) // 수신내용 확인

	respBody := resp.Body
	respStatus := resp.StatusCode

	json.NewDecoder(respBody).Decode(&statusInfo)
	fmt.Println(statusInfo)

	return &statusInfo, model.WebStatus{StatusCode: respStatus}
}
