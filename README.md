## Cloning the SonarBCD Repository

To get started with sonarbcd repository, follow these steps to clone it to your local machine:

1. **Clone the Repository**
   - Execute the following command to clone the SonarBCD repository:
     `git clone git@github.com:t1nfoil/sonarbcd.git`
     
2. **Navigate to the Cloned Repository**
   - Once the cloning process is complete, navigate into the cloned directory using the `cd` command:
     `cd sonarbcd`

3. **Run the pre-compiled binary for your OS version**
   - These are located in the binaries/ folder off of the repository root.
   - sonarbcd.exe -> Windows
   - sonarbcd_linux -> Linux
    
## sonarbcd Usage ##

### Program Flags ###

The program accepts several command-line flags to customize its behavior:

- **-inputcsv**: Specifies the input CSV file to convert. Default value is `bcd.csv`.

- **-outputdir**: Specifies the directory to output the generated files. Default is `./generated-labels`.

- **-checkcsv**: When set, performs basic checks on the CSV file for errors.

- **-zipname**: When set, the name of the zipfile to generate (without the .zip extension), in the output directory. Defaults to generated-labels

### Usage Example ###

```
# Convert custom CSV file
$ sonarbcd.exe -inputcsv=mydata.csv

# Output to a specific directory
$ sonarbcd.exe -outputdir=./output

# Perform basic checks on the CSV file
$ sonarbcd.exe -checkcsv
```

## CSV Field Parameters ##

   ### Data Field Formats ###

1. **company_name:** 
   - Format: Text, eg: "Sonar Software"

2. **discounts_and_bundles_url:** 
   - Format: URL, eg: https://www.sonar.software

3. **acp:**
   - Format: Boolean (true/false)
   - Notes: This is the Affordability Connectivity Program, use "Yes" or "No" if this package does or does not apply under ACP respectively.

4. **customer_support_url:** 
   - Format: URL, eg: https://www.sonar.software

5. **customer_support_phone:** 
   - Format: Phone Number, eg: 702-447-1247

6. **network_management_url:** 
   - Format: URL, eg: https://www.sonar.software

7. **privacy_policy_url:** 
   - Format: URL, eg: https://www.sonar.software

8. **fcc_id:** 
   - Format: Text

9. **data_service_id:** 
   - Format: Text, eg: "SONAR100"
   - Notes: "This is your internal data service id, this is combined with fix_or_mobile and the fcc_id to create the unique plan id"

10. **data_service_name:** 
    - Format: Text, eg: "MaxSpeed 100"

11. **fixed_or_mobile:** 
    - Format: Text, eg: "Fixed" or "Mobile"

12. **data_service_price:** 
    - Format: Price (e.g., $###.###), eg: $70.00
    - Notes: This is the regular service price after introductory period is done.

13. **billing_frequency_in_months:** 
    - Format: Integer (Number of months), eg: 1

14. **introductory_period_in_months:** 
    - Format: Integer (Number of months), eg: 6

15. **introductory_price_per_month:** 
    - Format: Price (e.g., $###.##), eg: $50.00

16. **contract_duration:** 
    - Format: Integer (Number of months), eg: 12

17. **contract_url:** 
    - Format: URL, eg: https://www.sonar.software

18. **early_termination_fee:** 
    - Format: Price (e.g., $###.###), eg: $100.00

19. **dl_speed_in_kbps:** 
    - Format: Integer, eg: 100000, interpreted as Kbps and will be converted to Mbps with 1 place decimal precision (eg: 1.5 Mbps not 1.50 Mbps)
    - Format: Decimal, eg: 1.5, interpreted as Mbps and will be converted to 1 place decimal precision.
    - Notes: Any decimals ending in .0 are converted to whole numbers (eg: 100.0 Mbps is displayed as 100 Mbps)

20. **ul_speed_in_kbps:** 
    - Format: Integer, eg: 100000, interpreted as Kbps and will be converted to Mbps with 1 place decimal precision (eg: 1.5 Mbps not 1.50 Mbps)
    - Format: Decimal, eg: 1.5, interpreted as Mbps and will be converted to 1 place decimal precision.
    - Notes: Any decimals ending in .0 are converted to whole numbers (eg: 100.0 Mbps is displayed as 100 Mbps)

21. **latency_in_ms:** 
    - Format: Integer (Milliseconds), eg: 120

22. **data_included_in_monthly_price:** 
    - Format: Integer (GB), eg: 1000

23. **overage_fee:** 
    - Format: Price (e.g., $###.###), eg: $5.00

24. **overage_data_amount:** 
    - Format: Integer (GB), eg: 5

## Data Validation Checks (using flag -checkcsv) ##

The application performs several checks on the provided data to ensure its integrity and compliance with expected formats. Below are the key validations conducted:

1, **Introductory Period and Price:**
    - Ensures that both introductory_period_in_months and introductory_price_per_month are set.
    - Verifies if introductory_period_in_months can be successfully cast to an integer.
    - Checks if introductory_price_per_month adheres to the format of [$]###.## (allowing an optional $ sign).
    - Validates that the length of introductory_price_per_month string is at most 8 characters.
    - Confirms that introductory_price_per_month can be cast to a 2 decimal precision float64.

2. **Data Service Price:**
    - Confirms that data_service_price is defined.
    - Checks if data_service_price follows the format of [$]###.### (allowing an optional $ sign).
    - Verifies that the length of data_service_price string is at most 8 characters.
    - Ensures that data_service_price can be cast to a 2 decimal precision float64.

3. **Download and Upload Speeds:**
   - Verifies if `dl_speed_in_kbps` and `ul_speed_in_kbps` are present in the data.
      - Ensures that `dl_speed_in_kbps` and `ul_speed_in_kbps` can be successfully parsed as integers.
      - Ensures that `dl_speed_in_kbps` and `ul_speed_in_kbps` can be successfully parsed as floating point.
      - Validates that the values are integers within the range of 0 to 10,000,000.
      - Validates that the values are within the range of 0.00 to 10000.

Any errors encountered during these checks will be printed to the console, along with the corresponding row information from the CSV file.


