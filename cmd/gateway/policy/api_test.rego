package vanti.bsp.ew.TestApi

import future.keywords.in

default allow := false

allow {
  input.method == "GetTest"
  input.token_valid
}

allow {
  input.method == "UpdateTest"
  input.token_valid
}