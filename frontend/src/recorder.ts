export type RecorderData = {
    readonly mediaRecorder: MediaRecorder;
    readonly chunks: BlobPart[];
}

export type RecordedAudio = {
    readonly audioUrl: string;
    readonly blob: Blob;
    readonly chunks: BlobPart[];
}

export async function startRecording(onDataAvailable: () => void) : Promise<RecorderData> {
    console.log("starting recording...");

    const stream : MediaStream = await navigator.mediaDevices.getUserMedia({
        audio: true, // constraints: audio-only
    });

    const recorderData: RecorderData = {
        mediaRecorder: new MediaRecorder(stream),
        chunks: [],
    }

    recorderData.mediaRecorder.ondataavailable = (e) => {
        console.log("recording...");
        recorderData.chunks.push(e.data);
        onDataAvailable();
    }

    recorderData.mediaRecorder.start(50);

    return recorderData;
}

export function stopRecording(recorderData: RecorderData) : RecordedAudio {
    console.log("Stopping recording...");

    recorderData.mediaRecorder.stop();

    console.log("Recoded MIME type: " + recorderData.mediaRecorder.mimeType);

    const blobData = new Blob(
        recorderData.chunks,
        {
            type: recorderData.mediaRecorder.mimeType
        });

    return {
        blob: blobData,
        audioUrl: window.URL.createObjectURL(blobData),
        chunks: recorderData.chunks, //copy raws so that blob continues working after discarding RecorderData
    }
}

export async function sendRecording(recordedAudio: RecordedAudio) : Promise<string> {
    console.log("Sending recording...");

    const formData = new FormData();
    formData.append('recordedAudio', recordedAudio.blob);

    const result : Response = await fetch('/api/upload', {
        method: 'POST',
        body: formData,
    });

    if (!result.ok) {
        throw new Error('Upload failed');
    }

    const sttRes : Response = await fetch('/api/speech-to-text', {
        method: 'GET'
    });

    if (!sttRes.ok) {
        throw new Error('STT failed');
    }

    return sttRes.text();
}