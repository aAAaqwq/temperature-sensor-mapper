---
title: Building and optimizing KubeEdge edge device management practice cases
authors:
  - "@aAAaqwq"
approvers:
  - 
creation-date: 2025-07-02
last-updated: 2025-07-02
status: Implementable 
---
# KubeEdge device management practice case optimization
---
  - [Project Background](#project-background)
    - [1. Develop a temperature sensor Mapper plug-in based on the latest version of Mapper-Framework](#1-develop-a-temperature-sensor-mapper-plug-in-based-on-the-latest-version-of-mapper-framework)
    - [2. Optimize existing IoT cases](#2-optimize-existing-iot-cases)
    - [Goal](#goal)
  - [Design Details](#design-details)
    - [1. Mapper Plugin Development and Testing](#1-mapper-plugin-development-and-testing)
    - [2. kubeedge-counter-demo Case Optimization](#2-kubeedge-counter-demo-case-optimization)
  - [Planning](#planning)
---

## Project Background

### 1. Develop a temperature sensor Mapper plug-in based on the latest version of Mapper-Framework
Based on the latest version of Mapper-Framework, develop a Modbus protocol Mapper device management plug-in, and use Modbus simulation software to simulate a temperature sensor. The Mapper plug-in needs to implement the management of the sensor, including device data collection, status reporting and other functions, and write and control the device through the RESTful API provided by Mapper. Finally, the implemented code needs to be merged into the KubeEdge Example repository, and a complete instruction document is attached to ensure that users can easily understand and apply the case.

### 2. Optimize existing IoT cases
For the original Device-IoT cases in the KubeEdge Example repository, select at least one of them for iterative optimization, including kubeedge-counter-demo, web-demo, etc. The optimization focuses on improving compatibility, performance, and ease of use to ensure that users can seamlessly use these cases in the latest version of KubeEdge. The optimized cases need to be merged into the KubeEdge Example repository and come with complete usage documentation.

## Goal

## Design Details

### 1. Mapper Plugin Development and Testing


### 2. kubeedge-counter-demo Case Optimization

## Planning
