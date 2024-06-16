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
    if (navigator.mediaDevices.getDisplayMedia) {
        localVideo.srcObject = await navigator.mediaDevices.getDisplayMedia({ video: true, audio: true });
    } else if (navigator.mediaDevices.getUserMedia) {
        localVideo.srcObject = await navigator.mediaDevices.getUserMedia({ video: true, audio: true });
    }
    localVideo.srcObject.getTracks().forEach(track => {
        localPC.addTrack(track, localVideo.srcObject);
    });

    localPC.oniceconnectionstatechange = function (e) {
    };

    localPC.onicecandidate = function (e) {
        if (e.candidate) {
            console.log("localPC candidate: " + JSON.stringify(e.candidate))
            remotePC.addIceCandidate(e.candidate);
        }
    };

    localPC.createOffer({
        offerToReceiveAudio: false,
        offerToReceiveVideo: false
    }).then(
        function (desc) {
            console.log("localPC description: " + JSON.stringify(desc))
            localPC.setLocalDescription(desc);

            remotePC.onicecandidate = function (e) {
                if (e.candidate) {
                    console.log("remotePC candidate: " + JSON.stringify(e.candidate))
                    localPC.addIceCandidate(e.candidate);
                }
            };

            remotePC.onaddstream = function (e) {
                remoteVideo.srcObject = e.stream;
            };

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
            console.log("remotePC description: " + JSON.stringify(desc))
            remotePC.setLocalDescription(desc);

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