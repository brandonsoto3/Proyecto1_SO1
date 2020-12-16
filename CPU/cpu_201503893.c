#include <linux/module.h>
#include <linux/init.h>
#include <linux/kernel.h>   
#include <linux/fs.h>
#include <linux/proc_fs.h>
#include <linux/seq_file.h> 
#include <asm/uaccess.h> 
#include <linux/hugetlb.h>
#include <linux/sched.h>        
#include <linux/sched/signal.h> 

//DEFINIENDO EL TAMAÑO DEL BUFFER
#define BUFSIZE  	150

//DEFINIENDO INFORMACIÓN ACERCA DEL MODULO
MODULE_LICENSE("GPL");
MODULE_DESCRIPTION("Modulo que escribe informacion acerca del CPU");
MODULE_AUTHOR("BRANDON JAVIER SOTO CASTAÑEDA 201503893");

struct sysinfo inf;

int contador = 0;
long total_memoria = 0;

static int write_file(struct seq_file * m, void *v) {	
    cont = 1;
    si_meminfo(&inf);
    total_memoria 	= (inf.totalram * 4)/1024;    
    seq_printf(m,"{\n");
    seq_printf(m,"      \"struct_lista_procesos\":[\n");
    dfs(&init_task, m, 0);
    seq_printf(m,"      ]\n");
    seq_printf(m,"}\n");
    return 0;
}

static int open(struct inode *inode, struct  file *file) {
  return single_open(file, write_file, NULL);
}

static struct file_operations operaciones =
{    
    .open = open,
    .read = seq_read
};


static int iniciar(void)
{
    proc_create("cpu_201503893", 0, NULL, &operaciones);
    printk(KERN_INFO "BRANDON JAVIER SOTO CASTANEDA\n");
    return 0;
}
 
static void salir(void)
{
    remove_proc_entry("cpu_201503893", NULL);
    printk(KERN_INFO "Diciembre 2020\n");
}
 
module_init(iniciar);
module_exit(salir); 