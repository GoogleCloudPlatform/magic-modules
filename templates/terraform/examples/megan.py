import os
import subprocess
import shutil
import re

os.mkdir("tmpdir")

errors = []


with open("update_list.txt") as f:
    files = f.read().split("\n")
# for old_name in os.listdir("."):
    for old_name in files:

        if old_name.endswith(".tf.erb"):
            print ("========================================================================")
            print ("------------"+old_name+"------------")
            print ("========================================================================")

            # make some changes
            with open(old_name, "r+") as tf:
                current = tf.read()
                found = re.findall("\<\%.*\%\>", current)
                if len(found) == 0:
                    errors.append(old_name+ "--> not found")

                replaced = re.sub("\<\%.*\%\>", "REPLACE_ME", current)
                tf.seek(0)
                tf.write(replaced)
                tf.truncate()

            new_name = old_name[:-4]

            os.rename(old_name, "tmpdir/"+new_name)

            os.chdir("tmpdir")

            process = subprocess.Popen(["terraform", "init"], stdout=subprocess.PIPE, stderr=subprocess.PIPE)
            process.communicate()

            process = subprocess.Popen(["terraform", "0.12upgrade", "-yes"])
            process.communicate()

            if os.path.exists("versions.tf"):
                os.remove("versions.tf")
            else:
                errors.append(old_name+"--> no versions")

            os.chdir("..")
            os.rename("tmpdir/"+new_name, old_name)

            # revert some changes
            with open(old_name, "r+") as tf:
                current = tf.read()

                for f in found:
                    current = re.sub("REPLACE_ME", f, current, count=1)

                tf.seek(0)
                tf.write(current)

for x in errors:
    print (x)

shutil.rmtree("tmpdir")