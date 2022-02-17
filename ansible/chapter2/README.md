# Password generation

I used this command to generate the SHA512 hash used in the "Create user" task.

```bash
pass=`pwgen --secure --capitalize --numerals --symbols 12 1`

echo $pass | mkpasswd --stdin --method=sha-512; echo $pass
```
