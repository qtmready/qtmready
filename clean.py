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
