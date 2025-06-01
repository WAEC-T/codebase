# Analysis README

This directory contains scripts and resources for processing and analyzing data from the Otii system.

## Workflow Overview

1. **Data Formatting**  
  The `format_data` script processes raw data directly from Otii, converting it into a structured format suitable for further analysis and visualization. The notebooks expect the processed data to be placed in a `data/` folder within this directory. A separate baseline notebook is provided for working independently with baseline data.

2. **Data Scheme Creation**  
  After formatting, we manually create a data scheme based on the described data. This scheme is further processed for visibility and optimization as each dataset is processed individually to make expansion easier and possible changes more manageable.

3. **Example Data**  
  An example of the processed data can be found in the `example-data` directory. This sample can be used with the `plot` and `stattest` scripts to generate plots and perform statistical analyses on the results.

## Directory Structure

- `example-data/` – Contains example processed data.

## Usage

1. Run the formatting scripts on raw Otii data.
2. Use the processed data and data scheme for further analysis.
3. Refer to `example-data` for sample input to plotting and statistical tools.
