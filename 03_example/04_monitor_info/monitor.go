// This example opens a window and closes it when the user presses the
// window's close button.
package main

import (
	"fmt"

	"github.com/phrasz/nag/allegro"
)

const fps int = 60

func main() {
	// ALLEGRO_MONITOR_INFO info;
	// int num_adapters;
	// int i, j;
	//
	// (void)argc;
	// (void)argv;
	//
	// if (!al_init()) {
	// 	 abort_example("Could not init Allegro.\n");
	// }

	// num_adapters = al_get_num_video_adapters();
	//
	// log_printf("%d adapters found...\n", num_adapters);
	//
	// for (i = 0; i < num_adapters; i++) {
	// 	 al_get_monitor_info(i, &info);
	// 	 log_printf("Adapter %d: ", i);
	// 	 log_printf("(%d, %d) - (%d, %d)\n", info.x1, info.y1, info.x2, info.y2);
	// 	 al_set_new_display_adapter(i);
	// 	 log_printf("   Available fullscreen display modes:\n");
	// 	 for (j = 0; j < al_get_num_display_modes(); j++) {
	// 			ALLEGRO_DISPLAY_MODE mode;
	//
	// 			al_get_display_mode(j, &mode);
	//
	// 			log_printf("   Mode %3d: %4d x %4d, %d Hz\n",
	// 				 j, mode.width, mode.height, mode.refresh_rate);
	// 	 }
	// }
	//

	allegro.Run(func() {
		// var (
		// 	display    *allegro.Display
		// 	eventQueue *allegro.EventQueue
		// 	running    = true
		// 	err        error
		// )

		numAdapters := allegro.NumVideoAdapters()

		fmt.Printf("%d adapters found...\n", numAdapters)
		for monNum := 0; monNum < numAdapters; monNum++ {
			monInfo, err := allegro.GetMonitorInfo(monNum)
			if err != nil {
				fmt.Printf("[ERROR] Cannot get Monitor Info! (%s)\n", err.Error())
				return
			}

			fmt.Printf("Adapter %d:\n", monNum)
			fmt.Printf("(%d, %d) - (%d, %d)\n", monInfo.X1(), monInfo.Y1(), monInfo.X2(), monInfo.Y2())

			allegro.SetNewDisplayAdapter(monNum)
			fmt.Printf("   Available fullscreen display modes:\n")

			numDisplayMode := allegro.NumDisplayModes()
			for modeNum := 0; modeNum < numDisplayMode; modeNum++ {
				mode, err2 := allegro.GetDisplayMode(modeNum)
				if err2 != nil {
					fmt.Printf("[ERROR] Cannot get Monitor Mode Info! (%s)\n", err2.Error())
					return
				}
				fmt.Printf("   Mode %3d: %4d x %4d, %d Hz\n", modeNum, mode.Width(), mode.Height(), mode.RefreshRate())
			}
		}
	})
}
