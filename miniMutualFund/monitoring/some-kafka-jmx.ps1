# In your project root (where the compose file is)
New-Item -ItemType Directory -Force -Path .\monitoring | Out-Null

# Download the JMX Exporter jar (adjust version if needed)
Invoke-WebRequest `
  -Uri https://repo1.maven.org/maven2/io/prometheus/jmx/jmx_prometheus_javaagent/0.20.0/jmx_prometheus_javaagent-0.20.0.jar `
  -OutFile .\monitoring\jmx_prometheus_javaagent.jar

# Save the JMX config
@"
startDelaySeconds: 0
lowercaseOutputName: true
lowercaseOutputLabelNames: true
rules:
  - pattern: 'kafka.server<type=(.+), name=(.+)PerSec\\w*, topic=(.+)><>Count'
    name: kafka_$1_$2_total
    labels: { topic: "$3" }
    type: counter
  - pattern: 'kafka.server<type=(.+), name=(.+)PerSec\\w*><>Count'
    name: kafka_$1_$2_total
    type: counter
"@ | Set-Content -Encoding UTF8 .\monitoring\kafka-jmx.yml
