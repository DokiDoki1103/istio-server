package config

// V1alpha3ApiVersion 使用的 istio api 版本
const V1alpha3ApiVersion = "networking.istio.io/v1alpha3"

// DestinationRule istio 目标规则
const DestinationRule = "DestinationRule"

const Prefix = "rbd-istio"

const TcpFlow = Prefix + "-flow-tcp"

const HttpFlow = Prefix + "-flow-http"

const Degrade = Prefix + "-degrade"
