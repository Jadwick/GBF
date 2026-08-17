# Goblin's Best Friend
*a [GNX](https://github.com/MovaFlow/GNX) patcher utility.*

Goblin's Best Friend (GBF) is a simple utility to install and manage GNX installations without using G3M. The goal was to make installing the modding framework GNX as **EASY** as possible.

GBF can:
- Install GNX on a vanilla game
- Update GNX if previously installed by Goblin's Best Friend

GBF makes a backup of the original `data.win` that it uses for updating or unisntalling GNX.

## Latest downloads:

## Instructions
- Download the latest release for your system (only Windows builds are tested)
- Put the executable in the same directory as the `data.win` file. (The same folder as `GoblinNest.exe`)
  - If you have Goblin Nest on Steam, you can find this through *Manage -> Browse local files*
 <img width="175" height="152" alt="Screenshot 2026-08-14 003507" src="https://github.com/user-attachments/assets/9c84fbcc-7e60-42ec-bd5d-5753660fcac5" />

- Run GBF
    - For Windows users, just double-click it
    - For Linux users, you may need to make it executable first
- At the main screen, press `1` then `<enter>` to select 'Install'
- Follow the on-screen prompts, everything else should be automated
- After a successful installation **DO NOT use G3M** to launch your game
    - Run GoblinNest.exe itself, or launch directly through Steam

<img width="267" height="149" alt="Screenshot 2026-08-14 003419" src="https://github.com/user-attachments/assets/83faadd4-25ec-4442-a200-cde3291d21be" />

GBF will make a `gbf_data` folder to store backups, config files, and other data in.

## Installing Mods
Mods need to be specifically for GNX. You can tell if they are because once unzipped, they will contain a manifest.json.

### Ensure you have the correct folder structure
When you unzip your mods they need to be placed in the correct location for GNX to read them.

Correct folder structure should look like:
```
|---GoblinNest.exe
|---data.win
|---GNX_mods/
    |---<mod folder>
        |-manifest.json
        |-<more mod stuff>
```
