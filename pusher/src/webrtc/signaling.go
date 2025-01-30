/*
 * @Author: gongluck
 * @Date: 2025-01-29 22:41:22
 * @Last Modified by: gongluck
 * @Last Modified time: 2025-01-30 01:11:33
 */

package webrtc

import (
	"encoding/json"
	"log"
)

type Offer struct {
	SDP  string `json:"sdp"`
	Type string `json:"type"`
}

type Candidate struct {
	Candidate        string `json:"candidate"`
	SDPMid           string `json:"sdpMid"`
	SDPMLineIndex    int    `json:"sdpMLineIndex"`
	UsernameFragment string `json:"usernameFragment"`
}

type SignalMessage struct {
	Offer     *Offer     `json:"offer,omitempty"`
	Candidate *Candidate `json:"candidate,omitempty"`
	RoomID    string     `json:"roomId"`
	UniqueID  string     `json:"uniqueId"`
}

func HandleSignaling(message []byte) {
	var signalMessage map[string]json.RawMessage
	if err := json.Unmarshal(message, &signalMessage); err != nil {
		log.Printf("Error unmarshalling message: %v\n", err)
		return
	}

	if offer, ok := signalMessage["offer"]; ok {
		var offerData Offer
		if err := json.Unmarshal(offer, &offerData); err != nil {
			log.Printf("Error unmarshalling offer: %v\n", err)
			return
		}
		handleOffer(&offerData)
	} else if candidate, ok := signalMessage["candidate"]; ok {
		var candidateData Candidate
		if err := json.Unmarshal(candidate, &candidateData); err != nil {
			log.Printf("Error unmarshalling candidate: %v\n", err)
			return
		}
		handleCandidate(&candidateData)
	}
}

func handleOffer(offer *Offer) {
	// 处理 Offer
	log.Printf("Handling offer: %+v\n", offer)

	CallCppFunction()
}

func handleCandidate(candidate *Candidate) {
	// 处理 Candidate
	log.Printf("Handling candidate: %+v\n", candidate)
}
