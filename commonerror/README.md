# Common Error

This `commonerror` package standardizes error reporting between grpc services.

# Usage

Initialize a new common error like so:

```
commonError := commonerror.New(commonerror.ErrCodeServer, commonerror.  ErrMsgServer)
```

Take care to always use error codes provided by the commonerror package.

**Note**: Providing code 0 will return `nil` as 0 is reserved for success

---

Convert inbuilt error to common error:

```
commonError := commonerror.Convert(err)
```
