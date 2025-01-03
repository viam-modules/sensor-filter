# VIAM SENSOR-FILTER MODULE

This is a [Viam module](https://docs.viam.com/extend/modular-resources/) that implements the generic service API in order to filter sensor data according to the user's requested conditions.



## Getting started

To use this module, follow these instructions to [add a module from the Viam Registry](https://docs.viam.com/modular-resources/configure/#add-a-module-from-the-viam-registry) and select the `viam:generic:sensor-filter` model.

This module implements the `DoCommand()` method of the [generic service API](https://docs.viam.com/services/generic/#api).




## Configure your `sensor-filter` service

> [!NOTE]  
> Before configuring your service, you must [create a robot](https://docs.viam.com/manage/fleet/robots/#add-a-new-robot).

Navigate to the **Config** tab of your robot’s page in [the Viam app](https://app.viam.com/). Click on the + button to the right of your robot part name. Then, click **Service**. Select the `generic` type, then select the `sensor-filter` model. Enter a name for your service and click **Create**.

### Example
In the example below, a sensor was already configured on the robot. To configure your own sensor, navigate to the **Config** tab of your robot’s page in [the Viam app](https://app.viam.com/). Click on the + button to the right of your robot part name. Then, click **Component**. Find sensor in the drop down list, then select your type of sensor and [configure it appropriately](https://docs.viam.com/components/sensor/#configuration).


#### Example configuration

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
        "value": 2,
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


### Description of Attributes

The following attributes are available to configure your sensor-filter module:


| Name                          | Type   | Inclusion       | Description                                                                                                                                                                                                                                                                                                                                                                                                                                                                                          |
| ----------------------------- | ------ | ------------ | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `sensor_name`                 | string | **Required** |Sensor name to be used as input.                                                                                                                                                                                                                                                                                                                                                                                                                                               |
| `reading`             | string | **Required** |Model to be used to extract faces before computing embedding. See [available extractors](#extractors-and-encoders-available).                                                                                                                                                                                                                                                                                                                                                                        |
| `conditions`        | list    | Optional   | Each condition contains an operator and a value. Possible operators include: "gt", "gte", "lt", "lte", "eq", and "neq".  They correspond to >, >=, <, <=, =, and !=, respectively. Values can be anything. Without any conditions, the service will always return "true"                                                                                                                                                                                                                            
