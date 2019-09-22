#!/bin/bash
#Create Producer FileSystem
cd /home/ubuntu
mkdir data
sudo chown ubuntu:ubuntu data
sudo mkfs.btrfs /dev/nvme1n1
sudo mount -t auto -v /dev/nvme1n1 /home/ubuntu/data

#Create Consumer FileSystem
cd /home/ubuntu
mkdir data
sudo chown ubuntu:ubuntu data
sudo mkfs.xfs /dev/nvme1n1
sudo mount -t auto -v /dev/nvme1n1 /home/ubuntu/data