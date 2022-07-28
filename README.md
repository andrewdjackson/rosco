# Sensor Reference

<table><tbody><tr><td class="has-text-align-left" data-align="left"><strong>Acronym</strong></td><td class="has-text-align-left" data-align="left"><strong>Sensor</strong></td></tr><tr><td class="has-text-align-left" data-align="left">ATS<br>80x05_intake_air_temp</td><td class="has-text-align-left" data-align="left">Air Temperature Sensor<br>Located under the air filter. <br>Measures the air temperature.</td></tr><tr><td class="has-text-align-left" data-align="left">CTS<br>80x03_coolant_temp</td><td class="has-text-align-left" data-align="left">Coolant Temperature Sensor<br>Located under the injection unit. <br>Measures the coolant (engine) temperature.</td></tr><tr><td class="has-text-align-left" data-align="left">CAS (CPS)<br>80x19_crankshaft_position_sensor</td><td class="has-text-align-left" data-align="left">Crank Angle Sensor (Crankshaft Position Sensor)<br>Located at the front / side of the engine.<br>Measures the position of the engine relative to TDC and is used to calculate the RPM.</td></tr><tr><td class="has-text-align-left" data-align="left">MAP<br>80x07_map_kpa</td><td class="has-text-align-left" data-align="left">Manifold Absolute Pressure Sensor<br>Located in the ECU, connects via vacuum pipes and fuel trap from the rear of the inlet manifold.<br>Measures the engine load</td></tr><tr><td class="has-text-align-left" data-align="left">TPS<br>80x09_throttle_pot</td><td class="has-text-align-left" data-align="left">Throttle Potentiometer Sensor<br>Located right of the air filter connected to throttle linkage.<br>Measures the throttle position</td></tr><tr><td class="has-text-align-left" data-align="left">IAC(V)<br>80x12_iac_position</td><td class="has-text-align-left" data-align="left">Idle Air Control (Valve)<br>Controls the throttle / stepper motor to adjust air : fuel ratio when idling and smooth transition when lifting off the throttle.</td></tr></tbody></table>

# MemsFCR Log File Format and Applied Calculations to Raw Data
<table><thead><tr><td>Column</td><td>Description</td><td>Calculation Applied</td></tr></thead><tbody><tr><td>#time</td><td>event timestamp hh:mm:ss.sss</td><td></td></tr><tr><td>80x01-02_engine-rpm</td><td>engine rpm</td><td></td></tr><tr><td>80x03_coolant_temp</td><td>temperature in degrees Celsius read from the Coolant Temperature Sensor (CTS). This sensor can be found under the injector unit. An essential value in the air:fuel ratio calculation</td><td>value - 55</td></tr><tr><td>80x04_ambient_temp</td><td>not used by the ECU, always returns 255</td><td>value - 55 = 200</td></tr><tr><td>80x05_intake_air_temp</td><td>temperature in degrees Celsius read from the Air Intake Temperature Sensor (ATS). This sensor can be found under the air filter. An essential value in the air:fuel ratio calculation</td><td>value - 55</td></tr><tr><td>80x06_fuel_temp</td><td>not used by the ECU, always returns 255</td><td>value - 55 = 200</td></tr><tr><td>80x07_map_kpa</td><td>manifold absolute pressure (MAP). Reads pressure from the back of the injector unit via the vacuum pipes and fuel trap. An essential value in the air:fuel ratio calculation</td><td></td></tr><tr><td>80x08_battery_voltage</td><td>the battery voltage. A figure &lt;12 volts will cause running issues</td><td>value / 10</td></tr><tr><td>80x09_throttle_pot</td><td>throttle potentiometer position. used by the ECU to determine throttle position when controlling idle speed</td><td>value * 0.02</td></tr><tr><td>80x0A_idle_switch</td><td>shows the state of the throttle switch, fitted on early vehicles. On systems without an actual throttle switch, the value shown indicates whether the MEMS ECU has calculated that the throttle is closed by using the throttle position sensor. If the switch shows 'ON' when the throttle is closed, then the vehicle will not idle correctly and the closed throttle position may need to be reset. This procedure is performed by fully depressing and releasing the accelerator pedal 5 times within 10 or fewer seconds of turning on the ignition and then waiting 20 seconds.</td><td></td></tr><tr><td>80x0B_uk1</td><td>unknown value</td><td></td></tr><tr><td>80x0C_park_neutral_switch</td><td>used on vehicles with an automatic gearbox</td><td>true / false</td></tr><tr><td>80x0D-0E_fault_codes</td><td>ECU fault codes:&lt;br&gt;Coolant temp sensor fault (Code 1)&lt;br&gt;Inlet air temp sensor fault (Code 2)&lt;br&gt;Fuel pump circuit fault (Code 10)&lt;br&gt;Throttle pot circuit fault (Code 16)</td><td></td></tr><tr><td>80x0F_idle_set_point</td><td>adjusts the idle rpm by the value shown. Adjusting idle speed will modify this value</td><td></td></tr><tr><td>80x10_idle_hot</td><td>the number of IACV steps from fully closed (0) which the ECU has learned as the correct position to maintain the target idle speed with a fully warmed up engine. If this value is outside the range 10 - 50 steps, then this is an indication of a possible fault condition or poor adjustment.</td><td>value - 35</td></tr><tr><td>80x11_uk2</td><td>unknown value</td><td></td></tr><tr><td>80x12_iac_position</td><td>Inlet Air Control valve (IACV) position (relates to expected Stepper Motor position)</td><td></td></tr><tr><td>80x13-14_idle_error</td><td>idle speed offset (also known as idle speed deviation)</td><td></td></tr><tr><td>80x15_ignition_advance_offset</td><td>adjustment to the ignition timing</td><td></td></tr><tr><td>80x16_ignition_advance</td><td>ignition advance, value of 128 = 0</td><td>(value / 2) - 24</td></tr><tr><td>80x17-18_coil_time</td><td>coil timing in ms</td><td>value * 0.002</td></tr><tr><td>80x19_crankshaft_position_sensor</td><td>position of the crankshaft from the position sensor (CPS)</td><td></td></tr><tr><td>80x1A_uk4</td><td>unknown value</td><td></td></tr><tr><td>80x1B_uk5</td><td>unknown value</td><td></td></tr><tr><td>7dx01_ignition_switch</td><td>status of the ignition switch</td><td>true / false</td></tr><tr><td>7dx02_throttle_angle</td><td>shows the position of the throttle disc obtained from the MEMS ECU using the throttle potentiometer. This value should change from a low value to a high value as the throttle pedal is depressed.</td><td>value * 6 / 10</td></tr><tr><td>7dx03_uk6</td><td>unknown value</td><td></td></tr><tr><td>7dx04_air_fuel_ratio</td><td>the current air:fuel ratio</td><td>value / 10</td></tr><tr><td>7dx05_dtc2</td><td>diagnostic trouble code - unknown codes</td><td></td></tr><tr><td>7dx06_lambda_voltage</td><td>the voltage read from the lambda sensor</td><td>value * 5</td></tr><tr><td>7dx07_lambda_sensor_frequency</td><td>not used by the ECU, value reads 255</td><td></td></tr><tr><td>7dx08_lambda_sensor_dutycycle</td><td>not used by the ECU, value reads 255</td><td></td></tr><tr><td>7dx09_lambda_sensor_status</td><td>ECU O2 circuit status, 1 active</td><td></td></tr><tr><td>7dx0A_closed_loop</td><td>ECU has entered closed loop and uses the lambda sensor for determining air:fuel ratio</td><td></td></tr><tr><td>7dx0B_long_term_fuel_trim</td><td>long term fuel trim (LTFT) displays ECU value to adjust fuelling. value of 128 = 0</td><td>value - 128</td></tr><tr><td>7dx0C_short_term_fuel_trim</td><td>short term fuel trim (STFT) displays ECU value to adjust fuelling</td><td></td></tr><tr><td>7dx0D_carbon_canister_dutycycle</td><td>not used by ECU, value reads 0</td><td></td></tr><tr><td>7dx0E_dtc3</td><td>diagnostic trouble code - unknown codes</td><td></td></tr><tr><td>7dx0F_idle_base_pos</td><td>the base value to offset idle position from</td><td></td></tr><tr><td>7dx10_uk7</td><td>unknown value</td><td></td></tr><tr><td>7dx11_dtc4</td><td>diagnostic trouble code - unknown codes</td><td></td></tr><tr><td>7dx12_ignition_advance2</td><td>ignition advance</td><td>value - 48</td></tr><tr><td>7dx13_idle_speed_offset</td><td>idle speed offset used to adjust idle speed</td><td></td></tr><tr><td>7dx14_idle_error2</td><td>idle error</td><td></td></tr><tr><td>7dx14-15_uk10</td><td>unknown value</td><td></td></tr><tr><td>7dx16_dtc5</td><td>diagnostic trouble code - unknown codes</td><td></td></tr><tr><td>7dx17_uk11</td><td>unknown value</td><td></td></tr><tr><td>7dx18_uk12</td><td>unknown value</td><td></td></tr><tr><td>7dx19_uk13</td><td>unknown value</td><td></td></tr><tr><td>7dx1A_uk14</td><td>unknown value</td><td></td></tr><tr><td>7dx1B_uk15</td><td>unknown value</td><td></td></tr><tr><td>7dx1C_uk16</td><td>unknown value</td><td></td></tr><tr><td>7dx1D_uk17</td><td>unknown value</td><td></td></tr><tr><td>7dx1E_uk18</td><td>unknown value</td><td></td></tr><tr><td>7dx1F_uk19</td><td>unknown value</td><td></td></tr><tr><td>0x7d_raw</td><td>hexadecimal response from the ECU for command 0x7D</td><td></td></tr><tr><td>0x80_raw</td><td>hexadecimal response from the ECU for command 0x80</td><td></td></tr><tr><td>engine_running</td><td>engine is running</td><td>true / false</td></tr><tr><td>warming</td><td>engine is warming up to operating temperature</td><td>true / false</td></tr><tr><td>at_operating_temp</td><td>engine is at operating temperature</td><td>true / false</td></tr><tr><td>engine_idle</td><td>engine is idle</td><td>true / false</td></tr><tr><td>idle_fault</td><td>hot or cold idle speed or idle offset is outside expected parameters</td><td>true / false</td></tr><tr><td>idle_speed_fault</td><td>cold idle speed is outside expected parameters</td><td>true / false</td></tr><tr><td>idle_error_fault</td><td>idle offset is outside expected parameters</td><td>true / false</td></tr><tr><td>idle_hot_fault</td><td>hot idle speed is outside expected parameters</td><td>true / false</td></tr><tr><td>cruising</td><td>rpm is stable but not idle; engine is cruising (differentiates from idle)</td><td>true / false</td></tr><tr><td>closed_loop</td><td>ECU is operating in closed loop (using lambda to determine air:fuel ratio)</td><td>true / false</td></tr><tr><td>closed_loop_expected</td><td>expecting the ECU to be in closed loop</td><td>true / false</td></tr><tr><td>closed_loop_fault</td><td>closed loop fault</td><td>true / false</td></tr><tr><td>throttle_active</td><td>the throttle pedal is depressed</td><td>true / false</td></tr><tr><td>map_fault</td><td>MAP readings is outside expected parameters</td><td>true / false</td></tr><tr><td>vacuum_fault</td><td>MAP and Air:Fuel ratio are outside expected parameters indicating a possible vacuum pipe fault</td><td>true / false</td></tr><tr><td>iac_fault</td><td>IAC position invalid if the idle offset exceeds the max error, yet the IAC Position remains at 0</td><td>true / false</td></tr><tr><td>iac_range_fault</td><td>IAC readings outside expected parameters</td><td>true / false</td></tr><tr><td>iac_jack_fault</td><td>high jack count indicating a possible problem with the stepper motor, throttle cable adjustment, or the throttle pot</td><td>true / false</td></tr><tr><td>o2_system_fault</td><td>detected a potential o2 system fault</td><td>true / false</td></tr><tr><td>lambda_range_fault</td><td>lambda sensor readings are outside expected parameters</td><td>true / false</td></tr><tr><td>lambda_oscillation_fault</td><td>lambda sensor not oscillating as expected</td><td>true / false</td></tr><tr><td>thermostat_fault</td><td>coolant temperature changes over time indicate thermostat fault (could also be a CPS fault)</td><td>true / false</td></tr><tr><td>crankshaft_sensor_fault</td><td>crankshaft position sensor (CPS) reading is outside expected parameters</td><td>true / false</td></tr><tr><td>coil_fault</td><td>coil is outside expected parameters</td><td>true / false</td></tr></tbody></table>

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

# Operational Status
## Is Engine Running?
```mermaid
flowchart LR
Engine("80x01-02_engine-rpm (RPM) > 0") --> EngineRunning(["engine_running = true"])
EngineRunning --> EngineStartTime(["record engine start time"])
```
## Is Engine at Operating Temperature?
```mermaid
flowchart LR
Temp("80x03_coolant_temp (CTS) > 80&deg;C") -- "false" ---> EngineWarming(["engine_warming = true"])
Temp -- "true" ---> EngineWarm(["engine_at_operating_temp = true"])
```
## Is Engine at Idle?
```mermaid
flowchart LR
Running(engine_running = true) --> Metric
Metric("7dx02_throttle_angle * 6 / 10 <= 14&deg") --> Result([engine_idle = true])
```
## Is Throttle Active?
```mermaid
flowchart LR
Engine("80x01-02_engine-rpm (RPM)> 1300") --> ThrottleActive(["throttle_active = true"])
Throttle("7dx02_throttle_angle * 6 / 10 > 14&deg;") --> ThrottleActive
```
## Is Closed Loop (O2 System) Active?
```mermaid
flowchart LR
Temp("7dx0A_closed_loop > 0") --> ClosedLoop(["closed_loop = true"])
```

# Operational Faults (Diagnosed)

## Is Battery Voltage too low?
```mermaid
flowchart LR
Metric("80x08_battery_voltage < 13V") --> Result(["battery_low = true"])
```
## Is Coil faulty?
```mermaid
flowchart LR
Running("engine_running = true") --> BatteryCheck
BatteryCheck("battery_low = false") --> Metric
Metric("80x17-18_coil_time * 0.0002 > 4ms") --> Result(["coil_fault = true"])
```
## Is MAP too high?
```mermaid
flowchart LR
Running("engine_idle = true") --> Metric
Metric("80x07_map_kpa > 45") --> Result(["map_fault = true"])
```
## Is O2 System active?
```mermaid
flowchart LR
Metric("7dx09_lambda_sensor_status > 0") --> Result(["o2_system_active = true"])
```
## Engine Idle fault?
```mermaid
flowchart LR
Running("engine_idle = true") --> EngineAtTemp{{"engine_at_operating_temp"}}
EngineAtTemp -- "true" ---> MetricHot
EngineAtTemp -- "false" ---> MetricCold
MetricHot("7dx0F_idle_base_pos > 55") --> Result(["idle_fault = true"])
MetricCold("7dx0F_idle_base_pos < 45") --> Result(["idle_fault = true"])
```
## Engine Hot Idle fault?
```mermaid
flowchart LR
Running("engine_idle = true") --> PreCheck
PreCheck("engine_at_operating_temp = true") --> Metric
Metric("7dx0F_idle_base_pos") --> |" > 55"|Result(["coil_fault = true"])
```
## Idle Air Control (IAC) Fault?
```mermaid
flowchart LR
```
## Vacuum Pipe Fault?
```mermaid
flowchart LR
```
## Are Lambda readings out of range?
```mermaid
flowchart LR
```
## Is Jack Count too high?
```mermaid
flowchart LR
```
## Is Crankshaft Sensor (CPS/CAS) faulty? 
```mermaid
flowchart LR
```
## Is Lambda Sensor faulty?
```mermaid
flowchart LR
```
## Is Lambda Sensor oscillating?
```mermaid
flowchart LR
```
## Is Idle Speed faulty?
```mermaid
flowchart LR
```

# Idle / Cruise Speed Diagnostic Tree
```mermaid
flowchart LR
Throttle{{7dx02_throttle_angle}} --> |"< 4&deg;"|ThrottleIdle([throttle_active = false])
Throttle --> |"> 4&deg;"|ThrottleIdle(["throttle_active = true"])

ThrottleIdle --> QueryIdleOffset{{80x13-14_idle_error > 100}}
QueryIdleOffset --> IdleOffsetFault(["idle_error_fault = true"])
```
```mermaid
graph LR  
Start{{engine_running = true}} --> Ready(Sample for 20s)
Ready -- Every Second --> Ready 
Ready --> QueryThrottle{{throttle_active}}
QueryThrottle --> |true|QueryCruiseIdle{{"RPM stddev < 5%"}}  
QueryCruiseIdle --> IsCruising(["✪ cruising = true"])
 
QueryThrottle --> |false|QueryTemp{{at_operating_temp}}  
QueryTemp --> |true|QueryIdleHot{{10 > 80x10_idle_hot < 50}}  
QueryIdleHot --> IdleHotFault(["✪ idle_hot_fault = true"])
QueryTemp --> |false|QueryIdleCold{{900 > RPM < 1200}}  
QueryIdleCold --> QueryWarming{{CTS increasing temp}}  
QueryWarming --> |Yes|EngineWarming(["warming = true"])
QueryWarming --> |No|CTSFault(["thermostat_fault = true"])
```

# Lambda / O2 Diagnostic Tree
```mermaid
graph LR  
Start{{engine_running = true}} --> Ready(Sample for 20s)
Ready -- Every Second --> Ready 
Ready --> WarmEngine{{at_operating_temp = true}} 
WarmEngine --> QueryLambdaStatus{{7dx09_lambda_sensor_status = 1}}
QueryLambdaStatus --> |No|LambdaFault(["o2_system_fault = true"])
WarmEngine --> QueryLambdaOsc{{7dx06_lambda_voltage oscillating?}}
QueryLambdaOsc --> |No|LambdaOscFault(["lambda_oscillation_fault = true"])
LambdaOscFault --> LambdaFault
WarmEngine --> QueryLambdaVoltage{{10mV > 7dx06_lambda_voltage < 900mV}}
QueryLambdaVoltage --> |No|LambdaRangeFault(["lambda_range_fault = true"])
LambdaRangeFault --> LambdaFault
```

# MAP Diagnostic Tree
```mermaid
graph LR  
Start{{engine_running}} --> |true| Ready(Sample for 20s)
Ready -- Every Second --> Ready  
Ready --> QueryMAP{{30 > 80x07_map_kpa < 60}}
QueryMAP --> MAPF(["map_fault = true"])
Start --> |false|QueryMAPOff{{80x07_map_kpa < 90}}
QueryMAPOff --> MAPF
```

# ECU Sensor Diagnostics
```mermaid
graph LR  
Start{{engine_running = true}} --> Ready(Sample for 20s)
Ready -- Every Second --> Ready  
Ready --> QueryCAS{{80x19_crankshaft_position_sensor = 0}}  
QueryCAS --> CASF(["crankshaft_sensor_fault = true"])  
Ready --> QueryCOIL{{80x17-18_coil_time < 4ms}}  
QueryCOIL --> COILF(["coil_fault = true"])
```


