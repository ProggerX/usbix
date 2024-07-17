# usbix
## USBix - declarative storage.
#### Use the power of Go and Nix Flakes to make your USB storage really beautiful
## Deadly simple flake
- ```nix {
    inputs.usbix.url = "github:ProggerX/usbix";
	outputs = { usbix, ... }: {
        defaultApp."x86_64-linux" = usbix.defaultApp."x86_64-linux"; # Change this to your architecture.
		usbix = [
			{
				device = "/dev/sde";
				type = "fat32-partitions"; # Just fat32 partitions. Only this is supported yet.
				partitions = [
					{
						label = "USBMUSIC";
						size = 30;
						units = "%";
					}
					{
						label = "USBFILES";
						size = 70;
						units = "%";
					}
				];
			}
		];
	};
}```
- change what you want
- just run ```nix run <path> -- <path>```, ```nix run . -- .``` for example
