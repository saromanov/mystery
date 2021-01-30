package api

type PutRequest struct {
	Namespace string `json:"namespace"`
	Value string `json:"value"`
}

type PutResponse struct {

}