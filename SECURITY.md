# Security

Once a user starts DEXbot and tells the Client to start listening for new transactions, they will be prompted to enter their private key for each wallet they have connected. This allows for the Client to autonomously sign and send the transaction back to the Cortex after verification and validation of the transaction integrity. It is crucial to note that there are important steps that DEXbot takes to ensure safety and security of the private key. 

First, DEXbot never stores, retrieves, sends, logs, or shares the private key in any capacity, at any point. The private key is completely local to your machine and is only used when signing transactions. Furthermore, after entering your private key, the Client completely wipes the byte slice containing the private key so that it only exists in the data structure that is used to sign transactions.  The Client has no persistence and when the application is stopped, this data is destroyed. To ensure maximum safety, make sure you are only using official releases of trusted operating systems and taking proper measures to secure your computer. If you want to take a deep dive into the source code, the entire Client application is open sourced and available on the official DEXbot GitHub repo.


Below are a few security guidelines to follow.

1. No administrator or representative of DEXbot will ever contact you for your seed phrase or private key. If this does happen, notify a member of DEXbot team on our Discord Server directly so we are aware of the situation. Additionally, no administration or  representative of DEXbot will ever DM you on Discord or any other communication channel for support. All support is carried out either in the public community Discord channel or via encrypted email. Every email will be signed with a GPG signature to ensure authenticity.  

2. No admin or representative will ever send you a direct link to download or upgrade the client application. It is very important to **ONLY** download DEXbot from the Official GitHub repo. To avoid malicious actors, **NEVER** use any files claiming to be DEXbot that are sent from any person, organization or 3rd Party. 

3. We have provided checksum values for every release of DEXbot. Make sure that you **ALWAYS CHECK THE CHECKSUM** before running any application. For security, verify that the checksum of the client application matches exactly to the version checksum on the official DEXbot Client GitHub repo before entering in any private key. If your version checksum does not match the official version checksum, do not use the code. This would mean that the code is not authentic and the security of your wallet or tokens can not be guaranteed.
