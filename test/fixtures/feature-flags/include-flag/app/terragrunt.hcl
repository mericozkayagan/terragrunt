include "root" {
  path   = "${find_in_parent_folders()}"
  # merge_strategy = "deep"
}

inputs = {
  string_feature_flag = feature.string_feature_flag.value
  int_feature_flag = feature.int_feature_flag.value
  bool_feature_flag = feature.bool_feature_flag.value
}
