public class Nested {

    public static void sizeOf(int[][] x) {
        int size = 0;

        for(int i = 0; i < x.length; i++) {
            for(int j = 0; j < x[i].length; j++) {
                size ++;            
            }
        }

        System.out.println(size);
    }


    public static void main(String[] args) {
        int[][] x = {{1,2,3},{2},{3,4},{5}};

        sizeOf(x);
        
    }
}