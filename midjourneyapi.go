package midjourneyapi

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
)

const host = "https://api.midjourneyapi.io/v2"

type ImagineMode string

const ImagineModeFast ImagineMode = "fast"
const ImagineModeTurbo ImagineMode = "turbo"

const StatusWaitingToStart = "waiting-to-start"
const StatusRunning = "running"

type ResultRequest struct {
	TaskId   string `json:"taskId"`
	Position int    `json:"position,omitempty"`
}

type ResultResponse struct {
	Status     string  `json:"status,omitempty"`
	Percentage float64 `json:"percentage,omitempty"`
}

type Client struct {
	apiKey string
}

func NewClient(apiKey string) *Client {
	return &Client{apiKey: apiKey}
}

type ImagineRequest struct {
	Prompt      string      `json:"prompt"`
	Mode        ImagineMode `json:"mode,omitempty"`
	CallbackURL string      `json:"callbackURL,omitempty"`
}

type ImagineResponse struct {
	TaskId string `json:"taskId"`
}

func getArrayFirstOrEmpty[T string | int](a []T) T {
	var r T
	if len(a) >= 1 {
		r = a[0]
	}
	return r
}

func (self *Client) Imagine(prompt string, mode ImagineMode, callbackURL ...string) (string, error) {
	var result ImagineResponse
	err := self.postJson("/imagine", ImagineRequest{
		Prompt:      prompt,
		Mode:        mode,
		CallbackURL: getArrayFirstOrEmpty(callbackURL),
	}, &result)

	return result.TaskId, err
}

type ImagineResultResponse struct {
	ResultResponse
	ImageURL string `json:"image_url,omitempty"`
}

func (self *Client) ImagineResult(taskId string, position ...int) (*ImagineResultResponse, error) {
	var result ImagineResultResponse
	if err := self.postJson("/result", ResultRequest{
		TaskId:   taskId,
		Position: getArrayFirstOrEmpty(position),
	}, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

type DescribeResponse struct {
	TaskId string `json:"task_id"`
}

func (self *Client) Describe(image io.Reader, callbackURL ...string) (string, error) {
	var body bytes.Buffer
	mw := multipart.NewWriter(&body)
	w, err := mw.CreateFormFile("image", "image.jpg")
	if err != nil {
		return "", err
	}

	if _, err := io.Copy(w, image); err != nil {
		return "", err
	}

	if len(callbackURL) >= 1 {
		if err := mw.WriteField("callbackURL", callbackURL[0]); err != nil {
			return "", err
		}
	}

	if err := mw.Close(); err != nil {
		return "", err
	}

	req, err := http.NewRequest(http.MethodPost, host+"/describe", &body)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", mw.FormDataContentType())
	req.Header.Set("Authorization", self.apiKey)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}

	var result DescribeResponse
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return "", err
	}

	return result.TaskId, nil
}

type DescribeResultResponse struct {
	ResultResponse
	Content []string `json:"content,omitempty"`
}

func (self *Client) DescribeResult(taskId string) (*DescribeResultResponse, error) {
	var result DescribeResultResponse
	if err := self.postJson("/result", ResultRequest{
		TaskId: taskId,
	}, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

type SeedRequest struct {
	TaskId      string `json:"taskId"`
	CallbackURL string `json:"callbackURL"`
}

type SeedResponse struct {
	TaskId string
}

func (self *Client) Seed(taskId string, callbackURL ...string) (string, error) {
	var result SeedResponse
	err := self.postJson("/seed", SeedRequest{
		TaskId:      taskId,
		CallbackURL: getArrayFirstOrEmpty(callbackURL),
	}, &result)

	return result.TaskId, err
}

type SeedResultResponse struct {
	ResultResponse
	Seed string `json:"seed,omitempty"`
}

func (self *Client) SeedResult(taskId string) (*SeedResultResponse, error) {
	var result SeedResultResponse
	if err := self.postJson("/result", ResultRequest{
		TaskId: taskId,
	}, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

type UpscaleRequest struct {
	TaskId      string `json:"taskId"`
	Position    int    `json:"position"`
	CallbackURL string `json:"callback_url,omitempty"`
}

type UpscaleResponse struct {
	ImageURL string `json:"imageURL"`
}

func (self *Client) Upscale(taskId string, position int, callbackURL ...string) (string, error) {
	var result UpscaleResponse
	err := self.postJson("/upscale", UpscaleRequest{
		TaskId:      taskId,
		Position:    position,
		CallbackURL: getArrayFirstOrEmpty(callbackURL),
	}, &result)

	return result.ImageURL, err
}

type FaceswapRequest struct {
	TargetImageURL string `json:"targetImageURL"`
	FaceImageURL   string `json:"faceImageURL"`
}

type FaceswapResponse struct {
	ImageURL string `json:"imageURL"`
}

func (self *Client) Faceswap(targetImageURL string, faceImageURL string) (string, error) {
	var result FaceswapResponse
	err := self.postJson("/faceswap", FaceswapRequest{
		TargetImageURL: targetImageURL,
		FaceImageURL:   faceImageURL,
	}, &result)

	return result.ImageURL, err
}

func (self *Client) postJson(path string, request interface{}, response interface{}) error {
	var body bytes.Buffer
	if err := json.NewEncoder(&body).Encode(request); err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, host+path, &body)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", self.apiKey)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return err
	}

	return nil
}
