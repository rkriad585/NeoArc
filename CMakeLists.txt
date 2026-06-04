cmake_minimum_required(VERSION 3.20)
project(neoarc
  VERSION 1.5.0
  DESCRIPTION "Cross-Platform Command Obfuscation & Alias Execution System"
  LANGUAGES NONE
)

# This project uses Go and Python, not C/C++.
# This CMakeLists.txt is provided for IDE integration and
# build tooling compatibility only.

# Custom target aliases for common workflows
add_custom_target(build
  COMMAND ${CMAKE_MAKE_PROGRAM} -f "${CMAKE_SOURCE_DIR}/Makefile" build
  WORKING_DIRECTORY "${CMAKE_SOURCE_DIR}"
  COMMENT "Build CLI for current platform"
)

add_custom_target(build-all
  COMMAND ${CMAKE_MAKE_PROGRAM} -f "${CMAKE_SOURCE_DIR}/Makefile" build-all
  WORKING_DIRECTORY "${CMAKE_SOURCE_DIR}"
  COMMENT "Cross-compile for all platforms"
)

add_custom_target(test
  COMMAND ${CMAKE_MAKE_PROGRAM} -f "${CMAKE_SOURCE_DIR}/Makefile" test
  WORKING_DIRECTORY "${CMAKE_SOURCE_DIR}"
  COMMENT "Run all tests"
)

add_custom_target(clean
  COMMAND ${CMAKE_MAKE_PROGRAM} -f "${CMAKE_SOURCE_DIR}/Makefile" clean
  WORKING_DIRECTORY "${CMAKE_SOURCE_DIR}"
  COMMENT "Remove build artifacts"
)

add_custom_target(docker
  COMMAND ${CMAKE_MAKE_PROGRAM} -f "${CMAKE_SOURCE_DIR}/Makefile" docker
  WORKING_DIRECTORY "${CMAKE_SOURCE_DIR}"
  COMMENT "Build Docker image for CLI"
)

add_custom_target(release
  COMMAND ${CMAKE_MAKE_PROGRAM} -f "${CMAKE_SOURCE_DIR}/Makefile" release
  WORKING_DIRECTORY "${CMAKE_SOURCE_DIR}"
  COMMENT "Tag and prepare a release"
)

message(STATUS "NeoArc ${PROJECT_VERSION}")
message(STATUS "Languages: Go (CLI), Python (Server)")
message(STATUS "Use 'make <target>' or 'cmake --build . --target <target>'")
