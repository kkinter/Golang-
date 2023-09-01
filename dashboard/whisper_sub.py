import os
import whisper
from whisper.utils import get_writer

def find_mp4_files(directory):
    mp4_files = []
    for root, _, files in os.walk(directory):
        for file in files:
            if file.endswith(".mp4"):
                file_path = os.path.join(root, file)
                mp4_files.append(file_path)
    return mp4_files

script_directory = os.path.dirname(os.path.abspath(__file__))
mp4_files = find_mp4_files(script_directory)

output_directory = script_directory  # mp4 파일이 위치한 폴더와 동일한 위치에 출력 폴더를 생성

cnt = 0

# writer
vtt_writer = get_writer("vtt", output_directory)
srt_writer = get_writer("srt", output_directory)
txt_writer = get_writer("txt", output_directory)

for file_path in mp4_files:
    print(file_path)
    model = whisper.load_model("tiny.en")
    result = model.transcribe(audio=file_path)
    cnt += 1

    # 파일 이름과 확장자 분리
    file_name, _ = os.path.splitext(file_path)
    
    # 결과를 파일로 저장
    vtt_file_path = file_name + ".vtt"
    srt_file_path = file_name + ".srt"
    txt_file_path = file_name + ".txt"

    vtt_writer(result, vtt_file_path)
    srt_writer(result, srt_file_path)
    txt_writer(result, txt_file_path)

    print(cnt, "  /  ", len(mp4_files))
    print("Created files:", vtt_file_path, srt_file_path, txt_file_path)
    print("=======================================================")