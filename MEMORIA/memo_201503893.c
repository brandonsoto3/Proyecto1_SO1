//INCLUIMOS LAS BIBLIOTECAS NECESARIAS
#include <linux/proc_fs.h>
#include <linux/seq_file.h> 
#include <asm/uaccess.h> 
#include <linux/hugetlb.h>
#include <linux/module.h>
#include <linux/init.h>
#include <linux/kernel.h>   
#include <linux/fs.h>

//DEFINIENDO EL TAMAÑO DEL BUFFER
#define BUFSIZE  	150

//DEFINIENDO INFORMACIÓN ACERCA DEL MODULO
MODULE_LICENSE("GPL");
MODULE_DESCRIPTION("Modulo que escribe informacion acerca de la memoria RAM");
MODULE_AUTHOR("BRANDON JAVIER SOTO CASTAÑEDA 201503893");

struct sysinfo inf;

static int write_file(struct seq_file * archivo, void *v) {	
    si_meminfo(&inf);
    long total_memoria 	= (inf.totalram * 4);
    long memoria_libre 	= (inf.freeram * 4 );
    long memoria_utilizada = total_memoria - memoria_libre;
    seq_printf(archivo, "{\n");
    seq_printf(archivo, "       \"struct_lista_ram\":[\n");
    seq_printf(archivo, "               {\n");
    seq_printf(archivo, "                   \"Total_de_memoria_RAM_del_servidor\":%lu,\n", total_memoria / 1024);
    seq_printf(archivo, "                   \"Total_de_memoria_RAM_consumida\":%lu,\n",memoria_utilizada / 1024);
    seq_printf(archivo, "                   \"Porcentaje_de_consumo_de_RAM \":%i\n", (memoria_utilizada * 100)/total_memoria);
    seq_printf(archivo, "               }");
    seq_printf(archivo, "       ]\n");
    seq_printf(archivo, "}\n");
    return 0;
}

static int abrir(struct inode *inode, struct  file *file) {
  return single_open(file, write_file, NULL);
}

static struct file_operations operaciones =
{    
    .open = abrir,
    .read = seq_read
};

static int iniciar(void)
{
    proc_create("memo_201503893", 0, NULL, &operaciones);
    printk(KERN_INFO "201503893\n");
    return 0;
}
 
static void salir(void)
{
    remove_proc_entry("memo_201503893", NULL);
    printk(KERN_INFO "Sistemas Operativos 1 \n");
}
 
module_init(iniciar);
module_exit(salir); 


