# Crafted with ❤ at Breu, Inc. <info@breu.io>, Copyright © 2024.
#
# Functional Source License, Version 1.1, Apache 2.0 Future License
#
# We hereby irrevocably grant you an additional license to use the Software under the Apache License, Version 2.0 that
# is effective on the second anniversary of the date we make the Software available. On or after that date, you may use
# the Software under the Apache License, Version 2.0, in which case the following will apply:
#
# Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with
# the License.
#
# You may obtain a copy of the License at
#
# http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on
# an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the
# specific language governing permissions and limitations under the License.

import os
import re

def remove_copyright_block(file_path):
  """
  Removes the copyright block from a file.

  Args:
      file_path (str): The path to the source file.
  """
  with open(file_path, 'r') as file:
    content = file.read()

  # Find the copyright block and remove it
  content = re.sub(r"^(// Crafted with.*?)\n// specific language governing permissions and limitations under the License.", "", content, flags=re.DOTALL | re.MULTILINE)

  with open(file_path, 'w') as file:
    file.write(content)

def process_files(directory, ignored=[]):
  """
  Processes files in a directory, removing the copyright block.

  Args:
      directory (str): The directory containing the source files.
      ignored (list): A list of file names to ignore.
  """
  if len(ignored) == 0:
    ignored = [".DS_Store", "README.md", "LICENSE", "CONTRIBUTING.md", "CHANGELOG.md", "CODE_OF_CONDUCT.md"]

  for root, _, files in os.walk(directory):
    for file in files:
      if file not in ignored:
        file_path = os.path.join(root, file)
        remove_copyright_block(file_path)

# Example usage:
# Remove copyright blocks from files in "cmd", "internal", and "deploy" directories
dirs = ["cmd", "internal", "deploy"]

for d in dirs:
  process_files(d)
