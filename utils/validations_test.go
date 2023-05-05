package utils

func TestValidateCreateUserRequest(t *testing.T){
  tests = []struct{
    name string
    req *siogeneric.AwCreateUserRequest
    expectedErr error
  }{
    {name: "valid request", req: &siogeneric.AwCreateUserRequest{}, expectedErr: nil}
  }
}
