# Htlc-demo
The Htlc-demo helps to realize the exchange of fabric and Ethereum assets

#bob,alice is an example alice is a sender, bob is a buyer 

1:Configure the environment and deploy contract with qs.sh, the same things on the ethereum

2:Useing ca-user.sh to create user's ca-identity

3:Create user's wallet with creat-asset.sh

4:Initialize user wallet with init.sh

5:Query user's asset information and create user's(alice) htlc with query-createHtlc.sh

6:Query created htlc's information queryHtlc-withdraw.sh, another user(bob) get preimage to create htlc on the ethereum, 
 they alice use preimage to get bob asset on the ethereum,the same things on the fabric for bob

7:Query user asset, the transaction is end
