# `sensor-filter` modular component

This module implements the [Viam generic service API](https://docs.viam.com/services/generic/#api) in a `viam:generic:sensor-filter` model.
With this model, you can filter sensor data according to user-specified conditions, allowing you to return `true` or `false` based on other sensor readings.

Navigate to the **CONFIGURE** tab of your machine's page.
Click the **+** button, select **Component or service**, then select the `generic / sensor-filter` model provided by the [`sensor-filter` module](https://app.viam.com/module/viam/sensor-filter).
Click **Add module**, enter a name for your service, and click **Create**.

## Configure your `sensor-filter` service

On the new component panel, copy and paste the following attribute template into your service's **Attributes** box:

```json
{
  "sensor_name": "<string>",
  "reading": "<string>",
  "conditions": [
    {
      "value": <any>,
      "operator": "<string>"
    }
  ]
}
```

### Attributes

The following attributes are available for `viam:generic:sensor-filter` services:

| Name          | Type   | Inclusion | Description |
| ------------- | ------ | --------- | ----------- |
| `sensor_name` | string | Required  | Sensor name to be used as input |
| `reading`     | string | Required  | The specific reading from the sensor to filter |
| `conditions`  | list   | Optional  | Each condition contains an operator and a value. Possible operators include: "gt", "gte", "lt", "lte", "eq", and "neq". They correspond to >, >=, <, <=, =, and !=, respectively. Values can be anything. Without any conditions, the service will always return `true` |

### Example Configuration

```json
{
  "sensor_name": "mySensor",
  "reading": "distance",
  "conditions": [
    {
      "value": 2,
      "operator": "gt"
    },
    {
      "value": 10,
      "operator": "lt"
    }
  ]
}
```

Full configuration with example sensor:

```json
{
  "components": [
    {
      "name": "mySensor",
      "namespace": "rdk",
      "type": "sensor",
      "model": "viam:ultrasonic:sensor",
      "attributes": {
        "board": "local",
        "echo_interrupt_pin": "13",
        "trigger_pin": "11"
      }
    }
  ]
  "services": [
    {
      "attributes": {
        "sensor_name": "mySensor",
        "reading": "distance"
        "conditions": [
          {
            "value": 2,
            "operator": "gt"
          }
        ]
      }
      "name": "SF-module",
      "type": "generic",
      "namespace": "rdk",
      "model": "viam:generic:sensor-filter"
    }
  ]
}
```