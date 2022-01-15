# Security

Upon downloading, installing, and running the application, you are required to enter your wallet private key during startup of the DEXbot application in order to use the application. The reason your private key is needed is that it is required to sign transactions. With that said, The client application does not transmit or store your private key at any time, and is explicitly used **only** to sign transactions for the automatic swaps. 

In order to maintain integrity, transparency, and accountability, we've publicly released the source code for the client application on on GitHub client application repo. By doing this, it allows you to comb through the source code files yourself to verify that your private key is never stored, shared, or transmitted, either locally, or on DEXbot's servers. 

NOTE: With this transparency, there comes the possibility that a malicious actor may download, edit, and distribute the altered codebase with intent to steal your private key. We've implemented security checks to the best of our ability to ensure the security of the application, but you must download the application from a verified source to ensure that it's a genuine version.

1. No administrator or representative of DEXbot will ever contact you to for your seed phrase or private key. If this does happen, notify a member of DEXbot team on our Discord Server directly so we are aware of the situation.

2. No admin or representative will ever send you a direct link to download or upgrade the client application. **The only secure methods to download or upgrade the client application are to either download directly from this GitHub repo, or directly from our website (DEXbot.io)**.

3. The client app includes information on using a checksum validator that you can compare to our GitHub version. For security, verify that the checksums of the client application and the checksum that is posted on our GitHub project or website matches exactly before entering in any private key. If your version checksum does not match ours, we cannot ensure the security of your wallet or tokens/coin contained within it since it will likely be an altered form of our codebase.

To view version changelogs and checksums, head over to our [Wiki](https://www.github.com/)
