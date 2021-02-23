public class Sample1 {
	public static void main(String[] args) {
		int v = 1;
		switch (v) {
			case 1: {
				String say = "yes";
				System.out.println("say "+say); 
				}
				break;
			case 2: {
				String say = "no";
				System.out.println("say "+ say); 
				}
				break;
			default: {
				String say = "???";
				System.out.println("say "+ say); 
				}
				break;
		}
	}
}
