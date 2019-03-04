#!/bin/bash
export $(cat .env | xargs)
./flexitracker
