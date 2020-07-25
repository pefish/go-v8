
## go-v8


## 编译V8

```shell script

cd v8

git checkout 6.9.427.19

./third_party/depot_tools/gn gen `pwd`/out/v8build '--args=   is_component_build=false   is_debug=false   libcpp_is_static=false   symbol_level=1   treat_warnings_as_errors=false   use_custom_libcxx=false   use_sysroot=false   v8_deprecation_warnings=false   v8_embedder_string="-v8worker2"   v8_enable_gdbjit=false   v8_enable_i18n_support=false   v8_enable_test_features=false   v8_experimental_extra_library_files=[]   v8_extra_library_files=[]   v8_imminent_deprecation_warnings=false   v8_monolithic=true   v8_static_library=false   v8_target_cpu="x64"   v8_untrusted_code_mitigations=false   v8_use_external_startup_data=false   v8_use_snapshot=true'

./third_party/depot_tools/ninja -v -C `pwd`/out/v8build v8_monolith

```


