# Version Checksums

A checksum is a technique used to determine the authenticity of data by running a value through a hashing algorithm. If you havent installed DEXbot yet, you can [download the client application or compile from the source code](https://github.com/DEXbotLLC/DEXbot_Client/wiki/2.-Installation).

a are established by the running a checksum validator on the binary file generated from the downloaded source code. After compiling a binary from the source code, or downloading the binary directly in [Releases](../../releases), open a terminal prompt, navigate to the folder containing the binary file, and run the following command to display the checksum of the DEXbot client application

<br>
<br>

**After downloading the application, it is very important to check the checksum of the binary before running it.** This ensures that the code you have compiled is safe and not has not been changed. Navigate to the folder with the compiled DEXbot application and run the following command.

<br>

| Operating System | Command |
| --------| ------------ |
|   Linux   | `sha256sum dexbot`   |
|   MacOS   | `openssl sha256 dexbot-mac`   |
|   Windows  | `certutil -hashfile dexbot-windows.exe SHA256`   |

<br>
<br>

Once you have the checksum value, you will need to verify that it is the correct value. Compare your value to the matching official version checksum below. 

<br>

## Release Date / Checksums

| Version |       OS         | Supported | Checksum                                                          |
| --------| ---------------- | ------------------- | ----------------------------------------------------------------- | 
| 0.2.1   |      Linux       |  :green_circle: |`1547aa615b05b54648138fbd41a8720314bd7ed4fc2b75436e3e2da62eb26bd8` | 
| 0.2.1   |      MacOS       |  :green_circle: |`5ffdc02959d4004fa2e26da646fbeb3d4dce96e23aeb103dedfc6d6f5444b3f5` | 
| 0.2.1   |      Windows     |  :green_circle: |`e549e23416740eac71c8d600af26ab43eb437e02d5e425d1b06d4236bfec5a8e` | 
| 0.2.0   |      Linux       |  :green_circle: |`bae871d9ed6cbba8ffe809de2420b392bad675e51a7d388dff84ee15b8a872cb` | 
| 0.2.0   |      MacOS       |  :green_circle: |`ff185565f7e981e58b1db5c07a63c245727ddfd20daa773db8d611af0c66faba` | 
| 0.2.0   |      Windows     |  :green_circle: |`ba58701991d9fde4a78a9946fe8bd9f286b7f202b4fb2dbe54a434b60a5913d0` | 
| 0.1.0   |      Linux       |  :x: |`9cf3f42ee56a29eeb18e057806e9f246dcf9117c62c4e5225f8d4be2fcac5944` | 
| 0.1.0   |      MacOS       |  :x: |`eb34e9bbf0f8739d5f65d8474721b395960dd4634a6a2888e4e8543b8a4ee661` | 
| 0.1.0   |      Windows     |  :x: |`4512e6c77c10c4840439dc3c944d2503ae03412998f8fdacfb53a0d6c51d3427` | 
