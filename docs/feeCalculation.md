## Fee Calculation Logic

### General Note

When a customer initiates an International Money Transfer the Payer Bank and the Beneficiary Bank may charge a fee for processing the payment. These fee charges may differ depending on the way the transfer has been requested.

When a customer makes an International Money Transfer he/she can choose who pays the transfer charges. A customer is usually offered one of the following three choices on who is to bear any charges in relation to the payment:

**"SHA" Transfer:** The Payer will pay fees to the Payer Bank i.e. Payer Bank's outgoing transfer charge. The Beneficiary will receive the amount transferred, minus the Beneficiary (intermediary) Bank charges. Unless requested otherwise, transfers are usually sent as "SHA".

**"OUR" Transfer:** All fees will be charged to the Payer - i.e. the Beneficiary gets the full amount sent. Any charges applied by the Beneficiary Bank will be billed back to the Payer (usually sometime after sending the payment).

**"BEN" Transfer:** The Payer does not pay any fee at all. The Beneficiary receives the sent payment minus all transfer charges (this includes both the Payer Bank and Beneficiary Bank fees).

In a SWIFT MT103, the following elements are of significance to the fee charges.


    * Field 71A - Details of Charges (BEN/OUR/SHA)
    * Field 71F - Sender's Charges
    * Field 71G - Receiver's Charges


### Chaincode Logic

The aim of the current chaincode logic is to accurately represent the movement of funds in the Statement Account of a given Financial Institution. Therefore it is imperative to understand how the above explained fee charges are embedded into the chaincode logic to articulate the final effect on the Statement Account.

#### Operating Model

<table>

<tbody>

<tr>

<th>Fee Type</th>

<th>Effect on the Payer Bank's Nostro Account</th>

</tr>

<tr>

<td>SHA</td>

<td>Upon Payment Confirmation, the Beneficiary Bank will debit the Statement Account for the amount stated in the Payment Instruction. The Beneficiary Bank will NOT debit it's transaction fee to the Statement Account as the Beneficiary will bear this cost. Hence beneficiary will receive the amount stated in the Payment Instruction minus the Beneficiary Bank fee.</td>

</tr>

<tr>

<td>OUR</td>

<td>Upon Payment Confirmation, the Beneficiary Bank will debit the Statement Account for the amount stated in the Payment Instruction and for the Beneficiary Bank's transaction fee. The Beneficiary will receive the full amount sent by the Payer.</td>

</tr>

<tr>

<td>BEN</td>

<td>The Payer Bank will create the Payment Instruction after deducting the Payer Bank's fee from the original amount received from the Payer. Therefore upon Payment Confirmation, the Beneficiary Bank will debit the Statement Account for the revised amount. The Beneficiary Bank will NOT debit it's transaction fee to the Statement Account as the Beneficiary will bear this cost. Hence beneficiary will receive the amount stated in the Payment Instruction minus the Beneficiary Bank fee.</td>

</tr>

</tbody>

</table>

#### Implementation Model

The fee logic uses the following fields:

*   **ReceiversCharge**: Indicates the fee that is owed by the Sending FI to the Receiving FI. Bulk settled at end of month.
*   **BenePays**: The amount the beneficiary will pay out of the transfer proceeds.
*   **Rebate**: In BEN / SHA arrangements, the amount owed by the Receiving FI to the Sending FI from the fee collected from the beneficiary.

The fee calculation feature implemented within the chaincode follows the below algorithm;

1.  Upon creating a Payment Instruction (`addPaymentInstruction`), the chaincode will calculate the `ReceiversCharge`, `BenePays` and `Rebate` using `calculateBankFees` function.

    The fees between financial institutions vary based on bilateral agreements. The following depicts the arrangements defined to date as per the [latest business requirements doc](https://github.com/ANZ-Blockchain-Lab/Corres-Bank-POC/blob/master/docs/PoC%20Business%20Requirements_V1.0.xlsx): ![](https://github.com/ANZ-Blockchain-Lab/Corres-Bank-POC/blob/master/docs/FeesAndRebatesLogic.png)

    **Note:** The rebate depicted in the above diagram ($12.50) was discontinued in 2016. This has been reflected in the chaincode.

    The Payment Instruction Struct will contain the Fee Type and the calculated fee amounts within the following variables respectively.

    *   `FeeType`
    *   `ReceiversCharge`
    *   `BenePays`  
    *   `Rebate`  

2.  When the Beneficiary Bank confirms the Payment Instruction, if the Fee Type is **OUR** the chaincode will call a function named `addLedgerEntryFee` to create the Ledger Entry on the Statement Account for the Beneficiary Bank fee amount. A Fee Ledger Entry will contain the Payment Instruction ID and the Payment Confirmation ID as references.
