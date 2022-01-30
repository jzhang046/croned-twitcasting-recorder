## **Croned Twicasting Recorder** 
Checks the live status of streamers on twitcasting.tv automatically at scheduled time, and records the live stream if it's available 

---

### **Disclaimer** 
This application constantly calls unofficial, non-documented twitcasting API to fetch live stream status. Please note that: 
* This application might not work in the future, subjecting to any change of twitcasting APIs 
* Checking live stream status at high frequency might result in being banned from using twitcasting service, subjecting to twitcasting's terms and condition

<span style="color:red">Please note the above and use this application at your own risk. </span>

---

### **Installation** 
* **Executables**   
  Executables can be found on [release page](https://github.com/jzhang046/croned-twitcasting-recorder/releases). 
* **Build from source**   
  Ensure that [golang is installed](https://golang.org/doc/install) on your system. 
  ```Bash
  git clone https://github.com/jzhang046/croned-twitcasting-recorder && cd croned-twitcasting-recorder
  go build -o ./bin/

  # Execute below command to start the recorder
  ./bin/croned-twitcasting-recorder
  ```


---

### **Configuration**
  Configuration file `config.yaml` must be present on the current directory. Please see [config_example.yaml] for example format.  
  At least 1 streamer should be specified in `config.yaml`  
  Multiple streamers could be specified with individual schedules. Status check and recording for different streamers would _not_ affect each other.  

  #### Field explanations: 
  + `screen-id`:  
    Presented on the URL of the screamer's top page.  
    Example: Top page URL of streamer [小野寺梓@真っ白なキャンバス](https://twitcasting.tv/azusa_shirokyan) is `https://twitcasting.tv/azusa_shirokyan`, the corresponding screen-id is `azusa_shirokyan`
  + `schedule`:   
    Please refer to the below docs for supported schedule definitions: 
    - https://pkg.go.dev/github.com/robfig/cron/v3#hdr-CRON_Expression_Format
    - https://pkg.go.dev/github.com/robfig/cron/v3#hdr-Predefined_schedules   

---

### **Output**  
  Output recording file would be put under the current directory, named after `screen-id-yyyyMMdd-HHmm.ts`  
  For example, a recording starts at 15:04 on 2nd Jan 2006 of streamer [小野寺梓@真っ白なキャンバス](https://twitcasting.tv/azusa_shirokyan) would create recording file `azusa_shirokyan-20060102-1504.ts`
