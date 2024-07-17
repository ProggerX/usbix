{
	description = "USBix - declare your file storage.";

	inputs = {
		nixpkgs.url = "github:nixos/nixpkgs/nixpkgs-unstable";
	};

	outputs = { self, nixpkgs, ... }:
	let systems = nixpkgs.lib.platforms.unix;
	in {
		defaultPackage = builtins.listToAttrs (builtins.map (system: { name = system; value =
		let pkgs = import nixpkgs { inherit system; };
		in pkgs.buildGoModule {
			name = "usbix";
			src = ./.;
			vendorHash = null;
		};}) systems );
		defaultApp = builtins.listToAttrs (builtins.map (system: { name = system; value =
		let pkgs = import nixpkgs { inherit system; };
		in {
			type = "app";
			program = let parted = "${pkgs.parted}/bin/parted"; usbix = "${self.defaultPackage.${system}}/bin/usbix"; in ''${pkgs.writeShellScript "usbix-run" "${parted} -v > /dev/null && sudo ${usbix} $@"}'';
		};}) systems );
	};
}
