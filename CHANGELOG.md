# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).
All dates in this file are given in the [UTC time zone](https://en.wikipedia.org/wiki/Coordinated_Universal_Time).

## Unreleased

## 0.1.2 - 2023-05-23

### Fixed

- When trying to determine the machine's serial number from a file (e.g. from the Raspberry Pi's firmware device tree), any trailing null terminators in the file are ignored.

## 0.1.1 - 2023-05-23

### Changed

- The machine name is now lazily loaded upon the first time the landing page is loaded and, unless an error occurs (e.g. from trying to determine the machine's serial number), cached for future use.

## 0.1.0 - 2023-05-20

### Added

- A minimalist static web page with a list of useful links.
