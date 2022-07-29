# Sensor Reference

<table><tbody><tr><td class="has-text-align-left" data-align="left"><strong>Acronym</strong></td><td class="has-text-align-left" data-align="left"><strong>Sensor</strong></td></tr><tr><td class="has-text-align-left" data-align="left">ATS<br>80x05_intake_air_temp</td><td class="has-text-align-left" data-align="left">Air Temperature Sensor<br>Located under the air filter. <br>Measures the air temperature.</td></tr><tr><td class="has-text-align-left" data-align="left">CTS<br>80x03_coolant_temp</td><td class="has-text-align-left" data-align="left">Coolant Temperature Sensor<br>Located under the injection unit. <br>Measures the coolant (engine) temperature.</td></tr><tr><td class="has-text-align-left" data-align="left">CAS (CPS)<br>80x19_crankshaft_position_sensor</td><td class="has-text-align-left" data-align="left">Crank Angle Sensor (Crankshaft Position Sensor)<br>Located at the front / side of the engine.<br>Measures the position of the engine relative to TDC and is used to calculate the RPM.</td></tr><tr><td class="has-text-align-left" data-align="left">MAP<br>80x07_map_kpa</td><td class="has-text-align-left" data-align="left">Manifold Absolute Pressure Sensor<br>Located in the ECU, connects via vacuum pipes and fuel trap from the rear of the inlet manifold.<br>Measures the engine load</td></tr><tr><td class="has-text-align-left" data-align="left">TPS<br>80x09_throttle_pot</td><td class="has-text-align-left" data-align="left">Throttle Potentiometer Sensor<br>Located right of the air filter connected to throttle linkage.<br>Measures the throttle position</td></tr><tr><td class="has-text-align-left" data-align="left">IAC(V)<br>80x12_iac_position</td><td class="has-text-align-left" data-align="left">Idle Air Control (Valve)<br>Controls the throttle / stepper motor to adjust air : fuel ratio when idling and smooth transition when lifting off the throttle.</td></tr></tbody></table>

# MemsFCR Log File Format and Applied Calculations to Raw Data
<table><thead><tr><td><strong>Metric</strong></td><td><strong>Description</strong></td><td><strong>Calculations Applied</strong></td></tr></thead><tbody><tr><td>#time</td><td>event timestamp hh:mm:ss.sss</td><td></td></tr><tr><td>80x01-02_engine-rpm</td><td>engine rpm</td><td></td></tr><tr><td>80x03_coolant_temp</td><td>temperature in degrees Celsius read from the Coolant Temperature Sensor (CTS). This sensor can be found under the injector unit. An essential value in the air:fuel ratio calculation</td><td>value - 55</td></tr><tr><td>80x04_ambient_temp</td><td>not used by the ECU, always returns 255</td><td>value - 55 = 200</td></tr><tr><td>80x05_intake_air_temp</td><td>temperature in degrees Celsius read from the Air Intake Temperature Sensor (ATS). This sensor can be found under the air filter. An essential value in the air:fuel ratio calculation</td><td>value - 55</td></tr><tr><td>80x06_fuel_temp</td><td>not used by the ECU, always returns 255</td><td>value - 55 = 200</td></tr><tr><td>80x07_map_kpa</td><td>manifold absolute pressure (MAP). Reads pressure from back of the injector unit via the vacuum pipes and fuel trap. An essential value in the air:fuel ratio calculation</td><td></td></tr><tr><td>80x08_battery_voltage</td><td>the battery voltage. A figure &lt;12 volts will cause running issues</td><td>value / 10</td></tr><tr><td>80x09_throttle_pot</td><td>throttle potentiometer position. used by the ECU do determine throttle position when controlling idle speed</td><td>value * 0.02</td></tr><tr><td>80x0A_idle_switch</td><td>shows the state of the throttle switch, fitted on early vehicles. On systems without an actual throttle switch the value shown indicates whether the MEMS ECU has calculated that the throttle is closed by using the throttle position sensor. If the switch shows 'ON' when the throttle is closed, then the vehicle will not idle correctly and the closed throttle position may need to be reset. This procedure is performed by fully depressing and releasing the accelerator pedal 5 times within 10 or less seconds of turning on the ignition and then waiting 20 seconds.</td><td></td></tr><tr><td>80x0B_uk1</td><td>unknown value</td><td></td></tr><tr><td>80x0C_park_neutral_switch</td><td>used on vehicles with an automatic gearbox</td><td>true / false</td></tr><tr><td>80x0D-0E_fault_codes</td><td>ECU fault codes:&lt;br&gt;Coolant temp sensor fault (Code 1)&lt;br&gt;Inlet air temp sensor fault (Code 2)&lt;br&gt;Fuel pump circuit fault (Code 10)&lt;br&gt;Throttle pot circuit fault (Code 16)</td><td></td></tr><tr><td>80x0F_idle_set_point</td><td>adjusts the idle rpm by the value shown. Adjusting idle speed will modify this value</td><td></td></tr><tr><td>80x10_idle_hot</td><td>the number of IACV steps from fully closed (0) which the ECU has learned as the correct position to maintain the target idle speed with a fully warmed up engine. If this value is outside the range 10 - 50 steps, then this is an indication of a possible fault condition or poor adjustment.</td><td>value - 35</td></tr><tr><td>80x11_uk2</td><td>unknown value</td><td></td></tr><tr><td>80x12_iac_position</td><td>Inlet Air Control valve (IACV) position (relates to expected Stepper Motor position)</td><td></td></tr><tr><td>80x13-14_idle_error</td><td>idle speed offset (also known as idle speed deviation)</td><td></td></tr><tr><td>80x15_ignition_advance_offset</td><td>adjustment to the ignition timing</td><td></td></tr><tr><td>80x16_ignition_advance</td><td>ignition advance, value of 128 = 0</td><td>(value / 2) - 24</td></tr><tr><td>80x17-18_coil_time</td><td>coil timing in ms</td><td>value * 0.002</td></tr><tr><td>80x19_crankshaft_position_sensor</td><td>position of the crankshaft from the position sensor (CPS)</td><td></td></tr><tr><td>80x1A_uk4</td><td>unknown value</td><td></td></tr><tr><td>80x1B_uk5</td><td>unknown value</td><td></td></tr><tr><td>7dx01_ignition_switch</td><td>status of the ignition switch</td><td>true / false</td></tr><tr><td>7dx02_throttle_angle</td><td>shows the position of the throttle disc obtained from the MEMS ECU using the throttle potentiometer. This value should change from a low value to a high value as the throttle pedal is depressed.</td><td>value * 6 / 10</td></tr><tr><td>7dx03_uk6</td><td>unknown value</td><td></td></tr><tr><td>7dx04_air_fuel_ratio</td><td>the current air:fuel ratio</td><td>value / 10</td></tr><tr><td>7dx05_dtc2</td><td>diagnostic trouble code - unknown codes</td><td></td></tr><tr><td>7dx06_lambda_voltage</td><td>the voltage read from the lambda sensor</td><td>value * 5</td></tr><tr><td>7dx07_lambda_sensor_frequency</td><td>not used by the ECU, value reads 255</td><td></td></tr><tr><td>7dx08_lambda_sensor_dutycycle</td><td>not used by the ECU, value reads 255</td><td></td></tr><tr><td>7dx09_lambda_sensor_status</td><td>ECU O2 circuit status, 1 active</td><td></td></tr><tr><td>7dx0A_closed_loop</td><td>ECU has entered closed loop and uses the lambda sensor for determining air:fuel ratio</td><td></td></tr><tr><td>7dx0B_long_term_fuel_trim</td><td>long term fuel trim (LTFT) displays ECU value to adjust fuelling. value of 128 = 0</td><td>value - 128</td></tr><tr><td>7dx0C_short_term_fuel_trim</td><td>short term fuel trim (STFT) displays ECU value to adjust fuelling</td><td></td></tr><tr><td>7dx0D_carbon_canister_dutycycle</td><td>not used by ECU, value reads 0</td><td></td></tr><tr><td>7dx0E_dtc3</td><td>diagnostic trouble code - unknown codes</td><td></td></tr><tr><td>7dx0F_idle_base_pos</td><td>the base value to offset idle position from</td><td></td></tr><tr><td>7dx10_uk7</td><td>unknown value</td><td></td></tr><tr><td>7dx11_dtc4</td><td>diagnostic trouble code - unknown codes</td><td></td></tr><tr><td>7dx12_ignition_advance2</td><td>ignition advance</td><td>value - 48</td></tr><tr><td>7dx13_idle_speed_offset</td><td>idle speed offset used to adjust idle speed</td><td></td></tr><tr><td>7dx14_idle_error2</td><td>idle error</td><td></td></tr><tr><td>7dx14-15_uk10</td><td>unknown value</td><td></td></tr><tr><td>7dx16_dtc5</td><td>diagnostic trouble code - unknown codes</td><td></td></tr><tr><td>7dx17_uk11</td><td>unknown value</td><td></td></tr><tr><td>7dx18_uk12</td><td>unknown value</td><td></td></tr><tr><td>7dx19_uk13</td><td>unknown value</td><td></td></tr><tr><td>7dx1A_uk14</td><td>unknown value</td><td></td></tr><tr><td>7dx1B_uk15</td><td>unknown value</td><td></td></tr><tr><td>7dx1C_uk16</td><td>unknown value</td><td></td></tr><tr><td>7dx1D_uk17</td><td>unknown value</td><td></td></tr><tr><td>7dx1E_uk18</td><td>unknown value</td><td></td></tr><tr><td>7dx1F_jack_count</td><td>unknown value</td><td></td></tr><tr><td>0x7d_raw</td><td>hexadecimal response from the ECU for command 0x7D</td><td></td></tr><tr><td>0x80_raw</td><td>hexadecimal response from the ECU for command 0x80</td><td></td></tr><tr><td>engine_running</td><td>engine is running</td><td>true / false</td></tr><tr><td>warming</td><td>engine is warming up to operating temperature</td><td>true / false</td></tr><tr><td>at_operating_temp</td><td>engine is at operating temperature</td><td>true / false</td></tr><tr><td>engine_idle</td><td>engine is idle</td><td>true / false</td></tr><tr><td>idle_fault</td><td>hot or cold idle speed or idle offset is outside expected parameters</td><td>true / false</td></tr><tr><td>idle_speed_fault</td><td>cold idle speed is outside expected parameters</td><td>true / false</td></tr><tr><td>idle_error_fault</td><td>idle offset is outside expected parameters</td><td>true / false</td></tr><tr><td>idle_hot_fault</td><td>hot idle speed is outside expected parameters</td><td>true / false</td></tr><tr><td>cruising</td><td>rpm is stable but not idle; engine is cruising (differentiates from idle)</td><td>true / false</td></tr><tr><td>closed_loop</td><td>ECU is operating in closed loop (using lambda to determine air:fuel ratio)</td><td>true / false</td></tr><tr><td>closed_loop_expected</td><td>expecting the ECU to be in closed loop</td><td>true / false</td></tr><tr><td>closed_loop_fault</td><td>closed loop fault</td><td>true / false</td></tr><tr><td>throttle_active</td><td>the throttle pedal is depressed</td><td>true / false</td></tr><tr><td>map_fault</td><td>MAP readings is outside expected parameters</td><td>true / false</td></tr><tr><td>vacuum_fault</td><td>MAP and Air:Fuel ratio are outside expected parameters indicating a possible vacuum pipe fault</td><td>true / false</td></tr><tr><td>iac_fault</td><td>IAC position invalid if the idle offset exceeds the max error, yet the IAC Position remains at 0</td><td>true / false</td></tr><tr><td>iac_range_fault</td><td>IAC readings outside expected parameters</td><td>true / false</td></tr><tr><td>iac_jack_fault</td><td>high jack count indicating possible problem with the stepper motor, throttle cable adjustment or the throttle pot</td><td>true / false</td></tr><tr><td>o2_system_fault</td><td>detected a potential o2 system fault</td><td>true / false</td></tr><tr><td>lambda_range_fault</td><td>lambda sensor readings are outside expected parameters</td><td>true / false</td></tr><tr><td>lambda_oscillation_fault</td><td>lambda sensor not oscillating as expected</td><td>true / false</td></tr><tr><td>thermostat_fault</td><td>coolant temperature changes over time indicate thermostat fault (could also be a CPS fault)</td><td>true / false</td></tr><tr><td>crankshaft_sensor_fault</td><td>crankshaft position sensor (CPS) reading is outside expected parameters</td><td>true / false</td></tr><tr><td>coil_fault</td><td>coil is outside expected parameters</td><td>true / false</td></tr></tbody></table>

# MemsFCR <-> ECU Initialisation Sequence
```mermaid
sequenceDiagram
autonumber
Note over MemsFCR, ECU: Initialisation
MemsFCR ->>+ ECU: Connect
ECU ->>- MemsFCR: Connected
MemsFCR ->>+ ECU: 0xCA
Note right of ECU: Initialise Command A 
ECU ->>- MemsFCR: 0xCA
MemsFCR ->>+ ECU: 0x75
Note right of ECU: Initialise Command B
ECU ->>- MemsFCR: 0x75 
MemsFCR ->>+ ECU: 0xF4
Note right of ECU: Heartbeat
ECU ->>- MemsFCR: 0xF4 
MemsFCR ->>+ ECU: 0xD0
Note right of ECU: ECU ID
ECU ->>- MemsFCR: 0x0D 0x99 0x99 0x99 0x99 
Note over MemsFCR, ECU: Request Dataframe Loop
loop Every second
activate MemsFCR
MemsFCR ->>+ ECU: 0x7D
Note right of ECU: Dataframe 0x7D
ECU ->>- MemsFCR: 0x7D XX XX XX.. XX
MemsFCR ->>+ ECU: 0x80
Note right of ECU: Dataframe 0x80
ECU ->>- MemsFCR: 0x80 XX XX XX.. XX
deactivate MemsFCR
end
```

# ECU Faults
```mermaid
graph LR
FaultCodesDTC0("80x0d-0e_fault_codes (0x0d byte)") --> |& b00000001|CTSFault(["CTS_fault = true"])
FaultCodesDTC0 --> |& b00000010|ATSFault(["ATS_fault = true"])
FaultCodesDTC1("80x0d-0e_fault_codes (0x0e byte)") --> |& b00000010|FuelPumpFault(["fuel_pump_circuit_fault = true"])
FaultCodesDTC1 --> |& b01000000|TPSFault(["TPS_circuit_fault = true"])
```
<br><br>
# Operational Status
## Is Engine Running?
If we have revs, the engine must be running :)<br>
The engine start time is recorded to determine when the operating temperature is expected to be reached and when the O2 system should activate.
```mermaid
flowchart LR
Engine("80x01-02_engine-rpm (RPM) > 0") --> EngineRunning(["engine_running = true"])
EngineRunning --> EngineStartTime(["record engine_start_time"])
```
## Is Engine at Operating Temperature?
Used in the diagnostics to determine if the engine has reached the operating temperature.
```mermaid
flowchart LR
Temp{{"80x03_coolant_temp (CTS) > 80&deg;C"}} -- "false" ---> EngineWarming(["engine_warming = true"])
Temp -- "true" ---> EngineWarm(["engine_at_operating_temp = true"])
```
## Is Engine at Idle?
RPM can't be used as a reliable indicator of idle so the throttle is used to determine whether the engine is at idle.
```mermaid
flowchart LR
Running(engine_running) --> Metric
Metric("7dx02_throttle_angle * 6 / 10 <= 14&deg") --> Result([engine_idle = true])
```
## Is Throttle Active?
If the RPM is above any cold start value OR the throttle is depressed (14&deg; is pretty conservative) then the throttle is active.
```mermaid
flowchart LR
Engine("80x01-02_engine-rpm (RPM)> 1300") --> ThrottleActive(["throttle_active = true"])
Throttle("7dx02_throttle_angle * 6 / 10 > 14&deg;") --> ThrottleActive
```
## Is Closed Loop Mode Active?
The ECU controls the air : fuel mix based on the readings from the primary sensors, lambda, MAP and temperature.<br>
If the ECU determines these systems are working correctly it will enter closed loop mode. The car is unlikely to pass the emissions tests if it cannot enter closed loop mode.
```mermaid
flowchart LR
Temp("7dx0A_closed_loop > 0") --> ClosedLoop(["closed_loop = true"])
```
<br><br>
# Operational Faults

## Is Battery Voltage too low?
The battery voltage is important for the sensors to provide correct values. A low battery can be the cause of a number of phantom faults.ß 
```mermaid
flowchart LR
Metric("80x08_battery_voltage < 13V") --> Result(["battery_low = true"])
```
## Is Coil faulty?
The time for the ignition coil to charge up to its specified current, as measured by the MEMS ECU. With a battery voltage of about 14V, this value should be about 2-3mS.<br> 
A high value for coil charge time may indicate a problem with the ignition coil primary circuit. A failing CAS sensor can also cause the coil time to increase.
```mermaid
flowchart LR
Running("engine_running") --> BatteryCheck
BatteryCheck("battery_low = false") --> Metric
Metric("80x17-18_coil_time * 0.0002 > 4ms") --> Result(["coil_fault = true"])
```
## Is MAP too high?
The MAP sensor should be under 45kPA (ideally 35kPA) when the engine is at idle. A high value indicates a vacuum fault.
```mermaid
flowchart LR
Running("engine_idle") --> Metric
Metric("80x07_map_kpa > 45") --> Result(["map_fault = true"])
```
## Engine Idle fault?
This Idle Base Position is the number of steps from 0 which the ECU will use as guide for starting idle speed control during engine warm up. <br>
The value will start at quite a high value (>100 steps) on a very cold engine and fall to < 50 steps on a fully warm engine.
A high value on a fully warm engine or a low value on a cold engine will cause poor idle speed control.<br>
Idle position is calculated by the ECU using the engine coolant temperature sensor.
```mermaid
flowchart LR
Running("engine_idle = true") --> EngineAtTemp{{"engine_at_operating_temp"}}
EngineAtTemp -- "true" ---> MetricHot
EngineAtTemp -- "false" ---> MetricCold
MetricHot("7dx0F_idle_base_pos > 55") --> Result(["idle_fault = true"])
MetricCold("7dx0F_idle_base_pos < 45") --> Result(["idle_fault = true"])
```
## Engine Hot Idle fault?
The Hot Idle is the number of IACV steps from fully closed (0) which the ECU has learned as the correct position to maintain the target idle speed with a fully warmed up engine.<br>
If this value is outside the range 10 - 50 steps, then this is an indication of a possible fault condition or poor adjustment. 
```mermaid
flowchart LR
Running("engine_idle") --> PreCheck
PreCheck("engine_at_operating_temp = true") --> EngineAtTemp
EngineAtTemp ---> MetricLow
EngineAtTemp ---> MetricHigh
MetricLow("80x10_idle_hot < 10") ---> Result(["idle_hot_fault = true"])
MetricHigh("80x10_idle_hot > 55") ---> Result(["idle_hot_fault = true"])
```
## Idle Air Control (IAC) Fault?
The stepper motor adjusts the idle air control valve (IACV) and is used to control engine idle speed and air flow from cold start up<br>
A high number of steps indicates that the ECU is attempting to close the stepper or reduce the airflow a low number would indicate the inability to increase airflow<br><br>
The position of the IACV stepper motor as calculated by the ECU. The ECU has no method of actually measuring this position but instead works it out by remembering how may steps it has moved the stepper since the last time the ignition was switched off. If a stepper motor fault exists, this number will be incorrect. This value will normally be changing during idle condition as the ECU makes minor changes to the idle speed. A value of 0 during idle conditions indicates a fault condition or poor adjustment, as does a very high value.
```mermaid
flowchart LR 
Running("engine_idle") --> PreCheck
PreCheck("7dx13_idle_speed_offset > 50") --> Metric
Metric("80x12_iac_position > 0") --> Result(["iac_fault = true"])
```
## Vacuum Pipe Fault?
The bane of all Mini Spi owners, the vacuum pipe system to read the manifold absolute pressure is fragile and prone to split or disconnected pipes.<br>
The MAP values are essential for the ECU to determine the load on the engine and have a significant impact on performance and even basic idling.<br>
When the engine is at idle, the MAP is expected to be less than 45kPa, in fact it should be around 35kPA. Higher than this will be due to a vacuum fault.
```mermaid
flowchart LR
Running("engine_idle") --> Metric
Metric("80x07_map_kpa > 45") --> Result(["vacuum_fault = true"])
```
## Is Jack Count too high?
On systems using a throttle body where the idle air is controlled by a stepper motor which directly acts on the throttle disk (normally metal inlet manifold), the count indicates the number of times the ECU has had to re-learn the relationship between the stepper position and the throttle position. If this count is high or increments each time the ignition is turned off, then there may be a problem with the stepper motor, throttle cable adjustment or the throttle pot. On systems using a plastic throttle body/manifold, the count is a warning that the MEMS ECU has never seen the throttle fully closed. The count is increased for each journey with no closed throttle, indicating a throttle adjustment problem.<br>
This "fault" is more of an indication that the jack count is higher than expected and is worth monitoring this value.
```mermaid
flowchart LR
Running("engine_idle") --> Metric
Metric("7dx1F_jack_count > 50") --> Result(["iac_jack_fault = true"])
```
## Is Crankshaft Sensor (CPS/CAS) faulty? 
The crankshaft sensor is used to determine the position of the engine so the ECU can manage the timing of the injection of fuel and ignition.<br>
A failing CAS will show as 0 entries or high spikes in the graph trace. When the ECU fails to get a CAS reading, it will suspend firing the ignitions (this will show as a high coil time)
```mermaid
flowchart LR
Metric("80x19_crankshaft_position_sensor is 0") --> Result(["crankshaft_sensor_fault = true"])
```
## Are Lambda readings out of range?
When the engine is running the lambda sensor voltage is expected to oscillate low to high rapidly.<br>
If the readings are above 900mV or below 10mV this indicates a failing lambda sensor.
```mermaid
flowchart LR
Running("engine_running") --> MetricLow
Running ---> MetricHigh
MetricLow("7dx06_lambda_voltage < 10") ---> Result(["lambda_range_fault = true"])
MetricHigh("7dx06_lambda_voltage > 900") ---> Result(["lambda_range_fault = true"])
```
## Is Lambda Sensor faulty?
When the engine is started the lambda sensor goes through a warm-up cycle. <br>
The lambda has a build in heater and the ECU activates this via a relay.<br>
The lambda should start oscillating from low to high volts after around 90s.<br>
A failure in this activity will cause the ECU to change the lambda status to 0 and run in open loop mode. This is a fault condition.
```mermaid
flowchart LR
Running("engine_running") --> PreCheck
PreCheck("time since engine_start_time > 90s") --> IsOscillating
subgraph Is Lambda Sensor oscillating?
    IsOscillating("standard deviation of 7dx06_lambda_voltage < 100") --> Result(["lambda_oscillation_fault = true"])
end
Result --> Metric
Metric{{"lambda_oscillation_fault is false"}} --> LambdaResult(["lambda_fault = true"])
```
## Is O2 System faulty?
```mermaid
flowchart LR
Metric("7dx09_lambda_sensor_status = 0") --> Result(["o2_system_fault = true"])
```

## Is Lambda Sensor oscillating?
The lambda sensor should oscillate low to high voltage at a high frequency.<br>
A sample is used for calculating the standard deviation, which is used to determine whether the lambda is working as expected.
```mermaid
flowchart LR
Running("engine_running") --> Metric
Metric("standard deviation of 7dx06_lambda_voltage < 100") --> Result(["lambda_oscillation_fault = true"])
```
## Is Thermostat faulty?
The coolant temperature is expected to rise at a rate of around 11s per degree. <br>
If the coolant temperature hasn't reached operating temperature (80&deg;C) by that time, this indicates the thermostat could be faulty.
```mermaid
flowchart LR
Running("engine_running") --> Step1
subgraph Determine if the coolant temperature has reached operating temp as expected
Step1("degrees_to_warm = 80&deg;C - 80x03_coolant_temp") --> Step2
Step2("expected_time_to_warm = 11s * degrees_to_warm") --> Step3
Step3("current_time > expected_time_to_warm") --> Step4
end
Step4("80x03_coolant_temp < 80&deg;C") --> Result(["thermostat_fault = true"])
```
## Is Idle Speed faulty?
Idle speed deviation indicates how far away the IAC is from the idle target<br>
Idle base position indicates the target for IAC position
A mean value of more than 100 indicates that the ECU is not in control of the idle speed and indicates a possible fault condition.
```mermaid
flowchart LR
Running("engine_running") --> Metric
Metric("mean of 7dx0F_idle_base_pos > 100") --> Result(["idle_speed_fault = true"])
```



