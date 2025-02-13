import json
import os
import sys
import re

def build_directory_structure(root_dir, fileExtensions):
    """
    Recursively builds a hierarchical JSON representation of the directory structure.
    """
    result = {"name": os.path.basename(root_dir), "type": "directory", "children": []}
    
    for entry in os.listdir(root_dir):
        full_path = os.path.join(root_dir, entry)
        
        if os.path.isdir(full_path):
            # If it's a directory, recurse into it
            # unless it is __pycache__ or .git
            if entry != "__pycache__" and entry != ".git":
                result["children"].append(build_directory_structure(full_path, fileExtensions))
        elif os.path.isfile(full_path):
            # Check if the file is a code file that is in the regex_matches list of regexes that we are interested in
            matchString = "\.("+ "|".join(fileExtensions) +")$"
            match = re.search(matchString, entry)
            if match:
                with open(full_path, "r", encoding="utf-8") as file:
                    content = file.read()
                result["children"].append({
                    "name": entry,
                    "type": "file",
                    "content": content
                })
    
    return result

def dirToStr(root_dir, fileExtensions):
    directory_structure = build_directory_structure(root_dir, fileExtensions)
    filesAndFolders = {
        "filesAndFolders": directory_structure
    }
    return json.dumps(filesAndFolders)

# get the json and format it with prety json
# print(json.dumps(json.loads(dirToStr("/home/user/Documents/GitHub/testBasedDev-demo-repo",["py","java","js","html","css","cpp","c","h","hpp","txt","md"])), indent=4))
# exit()


import pytest
import openai

if __name__ == "__main__":
    # save our current directory
    pwd = os.system("pwd")
    # get the path of the git repo from the env
    path = os.getenv("HOLDING_PATH")
    # run the tests
    exit_code = pytest.main(["--tb=short", "--disable-warnings", path])
    output = sys.stdout.getvalue()
    
    print(output)

    client = openai.OpenAI(api_key= os.getenv("OPENAI_API_KEY"), base_url= os.getenv("OPENAI_API_BASEURL"))

    model_engine = os.environ["MODEL"]

    prompt = os.getenv("PROMPT")
    with open(os.getenv("FILE"), "r") as f:
        prompt+= f"\n{f.read()}"
    print(f"Prompt: {prompt}")
    try:
        response = client.chat.completions.create(
            {
                'model': model_engine,
                'messages': [
                    {"role": "system", "content": "You are a coding assistant who works in the background. Your job is to take the provided test case and implement it."},
                    {"role": "user", "content": prompt},
                ]
            })
        if response.choices:
            review_text = response.choices[0].message.content.strip()
        else:
            review_text = f"No correct answer from OpenAI!\n{response.text}"
    except Exception as e:
        review_text = f"OpenAI failed to generate a review: {e}"

    print(f"Response>\n{review_text}")
    with open("review.txt", "w") as f:
        f.write(review_text)
