package smartcore.bos

import future.keywords.in

import data.scutil.rpc.read_request
import data.scutil.token.token_has_permission

# sc-bos private APIs that we want to treat as traits
traits := [
  "smartcore.bos.Access",
  "smartcore.bos.AnprCamera",
  "smartcore.bos.Button",
  "smartcore.bos.driver.dali.DaliApi",
  "smartcore.bos.EmergencyLight",
  "smartcore.bos.Meter",
  "smartcore.bos.MQTT",
  "smartcore.bos.Transport",
  "smartcore.bos.SecurityEvent",
  "smartcore.bos.ServiceTicket",
  "smartcore.bos.SoundSensor",
  "smartcore.bos.Status",
  "smartcore.bos.UDMI"
]

trait_request {
  some trait in traits
  startswith(input.service, trait)
}

allow {
  trait_request
  read_request
  token_has_permission("trait:read")
}

allow {
  trait_request
  token_has_permission("trait:write")
}

allow {
  trait_request
  token_has_permission("trait:*")
}