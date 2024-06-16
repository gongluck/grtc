/*
 * @Author: gongluck 
 * @Date: 2024-06-16 22:17:43 
 * @Last Modified by: gongluck
 * @Last Modified time: 2024-06-16 22:27:29
 */

'use strict';

const localVideo = document.getElementById('localVideo');
const remoteVideo = document.getElementById('remoteVideo');

const startPushBtn = document.getElementById('btnStartPush');
const startPullBtn = document.getElementById('bthStartPull');
const stopBtn = document.getElementById('bthStop');
startPullBtn.disabled = true;
stopBtn.disabled = true;

let localPC;
let remotePC;

startPushBtn.addEventListener('click', async function () {
    startPushBtn.disabled = true;
    startPullBtn.disabled = false;
    stopBtn.disabled = false;

    localPC = new RTCPeerConnection({});
    remotePC = new RTCPeerConnection({});

    if (navigator.mediaDevices.getDisplayMedia) {
        localVideo.srcObject = await navigator.mediaDevices.getDisplayMedia({ video: true, audio: true });
    } else if (navigator.mediaDevices.getUserMedia) {
        localVideo.srcObject = await navigator.mediaDevices.getUserMedia({ video: true, audio: true });
    }

    localVideo.srcObject.getTracks().forEach(track => {
        localPC.addTrack(track, localVideo.srcObject);
    });
    localPC.onicecandidate = function (e) {
        if (e.candidate) {
            console.log("localPC candidate: " + JSON.stringify(e.candidate));
            remotePC.addIceCandidate(e.candidate);
        }
    };

    remotePC.onaddstream = function (e) {
        remoteVideo.srcObject = e.stream;
    };
    remotePC.addEventListener('track', function (e) {
        if (remoteVideo.srcObject !== e.streams[0]) {
            remoteVideo.srcObject = e.streams[0];
        }
    });
    remotePC.onicecandidate = function (e) {
        if (e.candidate) {
            console.log("remotePC candidate: " + JSON.stringify(e.candidate))
            localPC.addIceCandidate(e.candidate);
        }
    };
});

startPullBtn.addEventListener('click', async function () {
    startPullBtn.disabled = true;
    
    localPC.createOffer().then(
        function (desc) {
            console.log("localPC description: " + JSON.stringify(desc))
            localPC.setLocalDescription(desc);
            // 交换sdp
            remotePC.setRemoteDescription(desc);
            remotePC.createAnswer().then(
                function (desc) {
                    console.log("remotePC description: " + JSON.stringify(desc))
                    remotePC.setLocalDescription(desc);

                    localPC.setRemoteDescription(desc);
                }
            );
        }
    );
});

stopBtn.addEventListener('click', async function () {
    startPushBtn.disabled = false;
    startPullBtn.disabled = true;
    stopBtn.disabled = true;

    localPC.close();
    localPC = null;
    remotePC.close();
    remotePC = null;

    localVideo.srcObject = null;
});
