# Golang CLI Application for managing multiply versions of any file. 

It must read .shadowrc from current directory and for each file entity listed in the config, it must create a shadow copy of the file/dir and store it in a hidden directory named .shadow.
The application should also provide commands to list all shadow copies, restore a shadow copy to its original location,
and delete a shadow copy. Multiple versions of the same file should be supported. Application should show all saved versions of the file.

Sample architecture: We always have 'virtual' copy, that actually HEAD in git terms, but we don't use git, we just use it as a reference.
We don't store actual file content in .shadow directory. this 'virtual' copy only and alias to current file state. 

Once user asked for creating a shadow copy, we create a copy of the file in a .shadow directory and store the metadata about the copy in a separate file. This metadata file will contain information such as the original file path, the timestamp of when the shadow copy was created, and any additional notes or tags provided by the user.

When listing shadow copies, we will read the metadata files to display relevant information about each shadow copy, such as the original file name, creation date, and any associated tags.

When user asks to restore a shadow copy, we sould ask user if she wants to save current state, and if yes, we should create a new shadow copy before restoring the selected one. Then we will copy the content of the shadow copy back to the original file location.

When user asks to delete a shadow copy, we will remove the corresponding metadata file and the shadow copy from the .shadow directory.

We are using https://github.com/charmbracelet stack for building pretty small minimalistic gui.


No work done on this yet.
