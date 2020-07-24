
#ifndef V8_WRAP_H
    #define V8_WRAP_H

    #ifdef __cplusplus  // 如果是c++
        extern "C" {
    #endif

    typedef void* IsolatePtr;

    void InitV8();
    void CloseV8();
    IsolatePtr NewIsolate();
    void CloseIsolate(IsolatePtr isolate);


    #ifdef __cplusplus
        }  // extern "C"
    #endif
#endif  // V8_WRAP_H

