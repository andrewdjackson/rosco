# Sensor Reference

<table><tbody><tr><td class="has-text-align-left" data-align="left"><strong>Acronym</strong></td><td class="has-text-align-left" data-align="left"><strong>Sensor</strong></td></tr><tr><td class="has-text-align-left" data-align="left">ATS<br>80x05_intake_air_temp</td><td class="has-text-align-left" data-align="left">Air Temperature Sensor<br>Located under the air filter. <br>Measures the air temperature.</td></tr><tr><td class="has-text-align-left" data-align="left">CTS<br>80x03_coolant_temp</td><td class="has-text-align-left" data-align="left">Coolant Temperature Sensor<br>Located under the injection unit. <br>Measures the coolant (engine) temperature.</td></tr><tr><td class="has-text-align-left" data-align="left">CAS (CPS)<br>80x19_crankshaft_position_sensor</td><td class="has-text-align-left" data-align="left">Crank Angle Sensor (Crankshaft Position Sensor)<br>Located at the front / side of the engine.<br>Measures the position of the engine relative to TDC and is used to calculate the RPM.</td></tr><tr><td class="has-text-align-left" data-align="left">MAP<br>80x07_map_kpa</td><td class="has-text-align-left" data-align="left">Manifold Absolute Pressure Sensor<br>Located in the ECU, connects via vacuum pipes and fuel trap from the rear of the inlet manifold.<br>Measures the engine load</td></tr><tr><td class="has-text-align-left" data-align="left">TPS<br>80x09_throttle_pot</td><td class="has-text-align-left" data-align="left">Throttle Potentiometer Sensor<br>Located right of the air filter connected to throttle linkage.<br>Measures the throttle position</td></tr><tr><td class="has-text-align-left" data-align="left">IAC(V)<br>80x12_iac_position</td><td class="has-text-align-left" data-align="left">Idle Air Control (Valve)<br>Controls the throttle / stepper motor to adjust air : fuel ratio when idling and smooth transition when lifting off the throttle.</td></tr></tbody></table>

# MemsFCR Log File Format and Applied Calculations to Raw Data
| Column | Description | Calculation Applied |
|--------|-------------|-------------|
| #time| event timestamp hh:mm:ss.sss ||
| 80x01-02_engine-rpm| engine rpm ||
| 80x03_coolant_temp| temperature in degrees Celsius read from the Coolant Temperature Sensor (CTS). This sensor can be found under the injector unit. An essential value in the air:fuel ratio calculation|value - 55|
| 80x04_ambient_temp| not used by the ECU, always returns 255| value - 55 = 200 |
| 80x05_intake_air_temp| temperature in degrees Celsius read from the Air Intake Temperature Sensor (ATS). This sensor can be found under the air filter. An essential value in the air:fuel ratio calculation|value - 55|
| 80x06_fuel_temp| not used by the ECU, always returns 255| value - 55 = 200 |
| 80x07_map_kpa| manifold absolute pressure (MAP). Reads pressure from back of the injector unit via the vacuum pipes and fuel trap. An essential value in the air:fuel ratio calculation | |
| 80x08_battery_voltage| the battery voltage. A figure <12 volts will cause running issues | value / 10 |
| 80x09_throttle_pot| throttle potentiometer position. used by the ECU do determine throttle position when controlling idle speed|value * 0.02
| 80x0A_idle_switch| shows the state of the throttle switch, fitted on early vehicles. On systems without an actual throttle switch the value shown indicates whether the MEMS ECU has calculated that the throttle is closed by using the throttle position sensor. If the switch shows 'ON' when the throttle is closed, then the vehicle will not idle correctly and the closed throttle position may need to be reset. This procedure is performed by fully depressing and releasing the accelerator pedal 5 times within 10 or less seconds of turning on the ignition and then waiting 20 seconds. ||
| 80x0B_uk1| unknown value ||
| 80x0C_park_neutral_switch| used on vehicles with an automatic gearbox | true / false |
| 80x0D-0E_fault_codes| ECU fault codes:<br>Coolant temp sensor fault (Code 1)<br>Inlet air temp sensor fault (Code 2)<br>Fuel pump circuit fault (Code 10)<br>Throttle pot circuit fault (Code 16) | |
| 80x0F_idle_set_point| adjusts the idle rpm by the value shown. Adjusting idle speed will modify this value | |
| 80x10_idle_hot| the number of IACV steps from fully closed (0) which the ECU has learned as the correct position to maintain the target idle speed with a fully warmed up engine. If this value is outside the range 10 - 50 steps, then this is an indication of a possible fault condition or poor adjustment. | value - 35 |
| 80x11_uk2| unknown value ||
| 80x12_iac_position| Inlet Air Control valve (IACV) position (relates to expected Stepper Motor position)  ||
| 80x13-14_idle_error| idle speed offset (also known as idle speed deviation)||
| 80x15_ignition_advance_offset| adjustment to the ignition timing ||
| 80x16_ignition_advance| ignition advance, value of 128 = 0 | (value / 2) - 24 |
| 80x17-18_coil_time| coil timing in ms | value * 0.002 |
| 80x19_crankshaft_position_sensor| position of the crankshaft from the position sensor (CPS) ||
| 80x1A_uk4| unknown value ||
| 80x1B_uk5| unknown value ||
| 7dx01_ignition_switch| status of the ignition switch | true / false |
| 7dx02_throttle_angle| shows the position of the throttle disc obtained from the MEMS ECU using the throttle potentiometer. This value should change from a low value to a high value as the throttle pedal is depressed. | value * 6 / 10 |
| 7dx03_uk6| unknown value ||
| 7dx04_air_fuel_ratio| the current air:fuel ratio | value / 10 |
| 7dx05_dtc2| diagnostic trouble code - unknown codes ||
| 7dx06_lambda_voltage| the voltage read from the lambda sensor | value * 5 |
| 7dx07_lambda_sensor_frequency| not used by the ECU, value reads 255 ||
| 7dx08_lambda_sensor_dutycycle| not used by the ECU, value reads 255 ||
| 7dx09_lambda_sensor_status| ECU O2 circuit status, 1 active  ||
| 7dx0A_closed_loop| ECU has entered closed loop and uses the lambda sensor for determining air:fuel ratio ||
| 7dx0B_long_term_fuel_trim| long term fuel trim (LTFT) displays ECU value to adjust fuelling. value of 128 = 0 | value - 128 |
| 7dx0C_short_term_fuel_trim| short term fuel trim (STFT) displays ECU value to adjust fuelling |  |
| 7dx0D_carbon_canister_dutycycle| not used by ECU, value reads 0 | |
| 7dx0E_dtc3| diagnostic trouble code - unknown codes ||
| 7dx0F_idle_base_pos| the base value to offset idle position from ||
| 7dx10_uk7| unknown value ||
| 7dx11_dtc4| diagnostic trouble code - unknown codes ||
| 7dx12_ignition_advance2| ignition advance | value - 48 |
| 7dx13_idle_speed_offset| idle speed offset used to adjust idle speed ||
| 7dx14_idle_error2| idle error ||
| 7dx14-15_uk10| unknown value ||
| 7dx16_dtc5| diagnostic trouble code - unknown codes ||
| 7dx17_uk11| unknown value ||
| 7dx18_uk12| unknown value ||
| 7dx19_uk13| unknown value ||
| 7dx1A_uk14| unknown value ||
| 7dx1B_uk15| unknown value ||
| 7dx1C_uk16| unknown value ||
| 7dx1D_uk17| unknown value ||
| 7dx1E_uk18| unknown value ||
| 7dx1F_uk19| unknown value ||
| 0x7d_raw| hexadecimal response from the ECU for command 0x7D |
| 0x80_raw| hexadecimal response from the ECU for command 0x80 |
| engine_running| engine is running | true / false |
| warming| engine is warming up to operating temperature | true / false |
| at_operating_temp| engine is at operating temperature | true / false |
| engine_idle| engine is idle | true / false |
| idle_fault| hot or cold idle speed or idle offset is outside expected parameters | true / false |
| idle_speed_fault| cold idle speed is outside expected parameters | true / false |
| idle_error_fault| idle offset is outside expected parameters | true / false |
| idle_hot_fault| hot idle speed is outside expected parameters | true / false |
| cruising| rpm is stable but not idle; engine is cruising (differentiates from idle) | true / false |
| closed_loop| ECU is operating in closed loop (using lambda to determine air:fuel ratio) | true / false |
| closed_loop_expected| expecting the ECU to be in closed loop | true / false |
| closed_loop_fault| closed loop fault | true / false |
| throttle_active| the throttle pedal is depressed | true / false |
| map_fault| MAP readings is outside expected parameters | true / false |
| vacuum_fault| MAP and Air:Fuel ratio are outside expected parameters indicating a possible vacuum pipe fault | true / false |
| iac_fault| IAC position invalid if the idle offset exceeds the max error, yet the IAC Position remains at 0  | true / false |
| iac_range_fault| IAC readings outside expected parameters | true / false |
| iac_jack_fault| high jack count indicating possible problem with the stepper motor, throttle cable adjustment or the throttle pot| true / false |
| o2_system_fault| detected a potential o2 system fault | true / false |
| lambda_range_fault| lambda sensor readings are outside expected parameters| true / false |
| lambda_oscillation_fault| lambda sensor not oscillating as expected | true / false |
| thermostat_fault| coolant temperature changes over time indicate thermostat fault (could also be a CPS fault) | true / false |
| crankshaft_sensor_fault| crankshaft position sensor (CPS) reading is outside expected parameters | true / false |
| coil_fault| coil is outside expected parameters | true / false |

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


