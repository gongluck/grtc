'use strict';

const localVideo = document.getElementById('localVideo');
const remoteVideo = document.getElementById('remoteVideo');

const startPushBtn = document.getElementById('btnStartPush');
const stopPushBtn = document.getElementById('bthStopPush');
const startPullBtn = document.getElementById('bthStartPull');
const stopPullBtn = document.getElementById('bthStopPull');

let localPC = new RTCPeerConnection({});
let remotePC = new RTCPeerConnection({});

startPushBtn.addEventListener('click', async function () {
    const stream = await navigator.mediaDevices.getDisplayMedia({ video: true, audio: true });
    localVideo.srcObject = stream;

    localPC.oniceconnectionstatechange = function (e) {
    };

    localPC.onicecandidate = function (e) {
        if (e.candidate) {
            console.log("localPC candidate: " + e.candidate.candidate)
            remotePC.addIceCandidate(e.candidate);
        }
    }

    localPC.addStream(stream);

    localPC.createOffer({
        offerToReceiveAudio: false,
        offerToReceiveVideo: false
    }).then(
        function (desc) {
            localPC.setLocalDescription(desc);

            // sdp交换
            remotePC.oniceconnectionstatechange = function (e) {
            }

            remotePC.onicecandidate = function (e) {
                if (e.candidate) {
                    console.log("remotePC candidate: " + e.candidate.candidate)
                    localPC.addIceCandidate(e.candidate);
                }
            }

            remotePC.onaddstream = function (e) {
                remoteVideo.srcObject = e.stream;
            }

            remotePC.setRemoteDescription(desc);
        }
    );
});

stopPushBtn.addEventListener('click', async function () {
    if (localPC) {
        localPC.close();
    }

    localVideo.srcObject = null;
});

startPullBtn.addEventListener('click', async function () {
    remotePC.createAnswer().then(
        function (desc) {
            remotePC.setLocalDescription(desc);

            // 交换sdp
            localPC.setRemoteDescription(desc);
        }
    );
});

stopPullBtn.addEventListener('click', async function () {
    if (remotePC) {
        remotePC.close();
    }

    remoteVideo.srcObject = null;
});