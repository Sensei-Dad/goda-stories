# Go-da Stories

Port of *Yoda Stories* to Go, based on [an article](https://www.gamedeveloper.com/programming/reverse-engineering-the-binary-data-format-for-star-wars-yoda-stories) about the datafile format.

To get started, simply copy the YODESK.DTA file from your Yoda Stories installation into `data/`.

## Goals

### Written in Go
First and foremost, this whole thing started as an excuse to learn Go. I also appreciate the purity that comes from using ONLY Go, despite how comfortable something might be to do with Python or Java or JS or any of the other languages that I already know. If I make a spaghetti monster of code, this time I'll do it in ONE language.

I shopped around for engines that would help with the "gamier" aspects of it, while still not doing everything for me. The design choices in [Ebiten](https://ebiten.org/) caught my eye: if "A dead simple 2D game library" isn't a perfect fit for something like *Yoda Stories*, I dunno what is. I could learn the language without learning an entire other game system with its *own* language on top of it. Having to wrap my poor little Object-oriented head around pointers and structs was bad enough. Ebiten hands me a game loop, throws in a graphics processor to draw stuff, and then gets out of my way. I love it.

Copying a [Roguelike tutorial](https://www.fatoldyeti.com/categories/roguelike-tutorial/) in Go to get started, I'm using [ByteArena's ECS library](https://github.com/bytearena/ecs) to hold the data, search through it, etc. A minute or two of searching through the Yoda Stories data file later, I'm convinced this is still the way to go.

### Not the Original, but Close Enough
To prevent myself going nuts and trying to code an entire cross-platform multiplayer nightmare that I'd never finish, I set some limitations before starting:

* **Limited to Go.** I'm learning Go, not trying to use some Eldritch necromancy with a decompiler to just import the original game's functions: I don't want to make a fancy interface for the existing game, I want to create a new one that uses the old game's assets. I'm also aware that 90% of this stuff would be already done, if I had used Godot: that's why I didn't ;-).
* **This is a programming project, not an art project.** Until the base game is finished, I'm allowed to use only the media assets extracted from the original game's data file. This applies to most everything in-game, because let's be fair: once that door opens I'll just be spending my time in GIMP mocking up Lightsaber "swooshes" from *Super Star Wars* everywhere instead of learning Go, because that pixel art is still fantastic. Let's keep it limited strictly to *Yoda Stories* assets...
    * ...with the exception of the user interface. Anything not using original game assets must be cobbled together "by hand" with Ebiten graphical tools and one grayscale PNG (originally, my own Dwarf Fortress 16x20 ASCII tileset), which I'll use to build the UI, dialogue, icons, controller buttons and other stuff. Unless I finish the base game and want to implement something like mods or new characters in later, I'll limit my graphical additions to this tileset only.
    * ...but I will, however, take advantage of Ebiten to tweak / flip / skew / attack our existing tiles with assorted Mathematics to get more out of them: *e.g.* skew the graphics for Luke's saber to be a wider slash, or implement crazier *even more modern* features that modern computers *might just* be able to pull off, like showing *more* than 9 rows of tiles at a time.
* **Write the game engine, not the game.** Related to the above. We're designing an engine that runs the original game, not trying to completely re-balance all of *Yoda Stories*. Use the original maps, scripting, and game logic wherever **feasible**:
    * Maps, Tiles
    * Items, Weapons, and enemies
    * Mission text: All the original Endings and scripts
    * Sounds & Music

Long story short: this doesn't need to be 100% frame-perfect-accurate. I don't want to write a mountain of code to fill in a tiny hole in the game logic. If I can't figure out what something does, it's okay to use the best guess at the time and move on, especially if it's considerably less work.

### But not THAT Close
When writing down all the things that I still *wanted* to adjust or change about the game, I realized that most of what I didn't like about the original (along with *Indy's Desktop Adventures*) was because of its Windows-3.1-native interface. And since native UI in Go at the time of this writing isn't universally "a thing", I'm sticking to using Ebiten to handle interpreting all the user input while giving the game a bit of a face-lift. Other gripes and notes on the interface design are in [their own file](/docs/Interface.md).

## Misc. Notes
Things I've recorded about my exploration of the game's data and my attempts to make an engine for it are all in the `docs/` folder:

* [/docs/Extraction.md] - Notes about the Yoda Stories datafile format
* [/docs/Interface.md] - Things related to the UI

## References
Here's a list of those who have pulled this off before me, or whose efforts got me pointed along the path that worked. A large part of what I've done is because these people laid such complete foundations for me to stand on:

* Of course, [the original game](https://en.wikipedia.org/wiki/Star_Wars:_Yoda_Stories) deserves a nod of respect. It's a gem. This and *Indy's Desktop Adventures* have enough about them to fit nicely into the "rogue-lite" genre, before it was even a thing.
* [The original article](https://www.gamedeveloper.com/programming/reverse-engineering-the-binary-data-format-for-star-wars-yoda-stories) about extracting the datafile format, by Zach Barth. (Not to mention the joys that he's brought with Infinifactory and SpaceChem, and all the other things inspired by his work.) I started this whole mess by trying to duplicate his code to view the tiles and make maps out of them, while learning to code in Go.
* After Go itself, [Ebiten](https://ebiten.org/) is the tool that makes everything else work. Hajime Hoshi has put together a wonderful bit of code. I spent some time shopping between Godot and Ebiten... and I'm glad I went with it. If Godot is a kitchen appliance that has a million-and-one functions plus a "popcorn" button, then it has that many more things to look up in the manual; Ebiten is the plain 10-inch chef's knife that looks unassuming but does its job flawlessly, and ends up being the most useful thing in your kitchen.
* [DesktopAdventures](https://github.com/shinyquagsire23/DesktopAdventures/blob/master/scrdoc.txt) by shinyquagsire23 is a re-implementation that was key in helping me decode what the IACT sections were doing.
* [WebFun](https://github.com/cyco/WebFun) by Cyco is another fantastic re-implementation of the original game engine, and the code there filled in all the gaps in my knowledge of the datafile that the above one didn't.
