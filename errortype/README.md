# ErrorType

This package provides a custom Error that allows one to:

- Differentiate between errors within and across packages in a project by their **code** and **package**.
- Obtain the **stack trace error** through the error message.

# Usage

Define ErrorType and use it to initialize new Error instances

```
import github.com/twothicc/common-go/errortype

var errorType1 = errortype.ErrorType{Code: 1, Pkg: "package1"}

...

func try() error {
    return errorType1.New("sample")
}
```

---

Wrap with a different errortype

```
error1 := errortype1.New("one")
error2 := errortype2.Wrap(error1)
error3 := errortype3.WrapWithMsg(error2, "three")
fmt.Println(error3.Error())
```

You should expect to see `error: code=3, pkg=package3, msg=three | error: code=2, pkg=package2, msg=one | error: code=1, pkg=package1, msg=one"`

---

Check if an error is of an ErrorType

```
errorType1.Is(error1)
```

You should expect to output true
