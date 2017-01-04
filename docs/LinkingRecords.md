# Linking Real World Records to Blockchain Records

![linkingrealworldrecords](https://cloud.githubusercontent.com/assets/15486303/17009296/ae8df636-4f3c-11e6-8af2-b16b824b54bd.png)

This is a high-level depiction of how we eliminate the need for matching two separately held records - by effectively linking them together via a shared DL.
  - When we create a request entry on the DL, we need to link it to a real world record (i.e. a Nostro Mirror entry). We do this through a method of calculating a `lookup key`.
    - All records in the DL are stored against a `lookup key`.
    - We generate the `lookup key` for the request by hashing the `MT103` message*.  
    - What this does is allow anyone with access to that `MT103` and the appropriate hashing function to generate the same `lookup key` and find the DL record that relates to the real world `MT103`.
  - In the same way that it allows the `Sending FI` to locate the DL record that corresponds to a real world MT103, it also allows the `Receiving FI` to do the same.
    - **This is where the innovation is**. By allowing the `Receiving FI` to locate a `request entry` and link its `confirmation entry` to it, we are creating a linkage between the `Sending FIs` Nostro Mirror record and the `Receiving FIs` Statement Account [MT103 <--> Request on DL <--> Confirmation on DL <--> Statement Account Transaction]. 
    - This link:
      - Communicates the completion of a transaction in real time, as opposed to the end of day batch timeframes;
      - Removes the need to match transactions based on their attributes, since links are established on creation of the records. This means the traditional matching rules are not required.

*Note: While we currently use a `MT103` to generate the `lookup key`, we could actually use any combination of information, so long as the following properties hold true:
  - The `Receiving FI` has access to the information that the `Sending FI` used to generate the `lookup key`;
  - The `Receiving FI` knows the exact combination and sequence of data to feed into the hash function to generate that `lookup key`. 
MT103s satisfy both criteria because both `Sending FI` and `Receiving FI` have access to that message via SWIFT. We are currently trying to identify a subset of the `MT103` data that can be used for any transaction type.

I encourage discussion on this and welcome any questions. This concept is a critical part of this POC, as it removes the need to replicate existing matching processes.
